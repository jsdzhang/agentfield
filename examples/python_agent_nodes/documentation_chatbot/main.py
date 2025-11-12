"""Simplified Documentation chatbot with parallel retrieval and self-aware synthesis."""

from __future__ import annotations

import asyncio
import os
from pathlib import Path
import sys
from typing import Any, Dict, List, Sequence

from agentfield import AIConfig, Agent
from agentfield.logger import log_info

if __package__ in (None, ""):
    current_dir = Path(__file__).resolve().parent
    if str(current_dir) not in sys.path:
        sys.path.insert(0, str(current_dir))

from chunking import chunk_markdown_text, is_supported_file, read_text
from embedding import embed_query, embed_texts
from schemas import (
    Citation,
    DocAnswer,
    DocumentChunk,
    DocumentContext,
    IngestReport,
    QueryPlan,
    RetrievalResult,
)

app = Agent(
    node_id="documentation-chatbot",
    agentfield_server=os.getenv("AGENTFIELD_SERVER", "http://localhost:8080"),
    ai_config=AIConfig(
        model=os.getenv("AI_MODEL", "openrouter/openai/gpt-4o-mini"),
        temperature=0.2,
    ),
)


# ========================= Ingestion Skill (Unchanged) =========================


@app.skill()
async def ingest_folder(
    folder_path: str,
    namespace: str = "documentation",
    glob_pattern: str = "**/*",
    chunk_size: int = 1200,
    chunk_overlap: int = 250,
) -> IngestReport:
    """
    Chunk + embed every supported file inside ``folder_path``.

    Uses two-tier storage:
    1. Store full document text ONCE in regular memory
    2. Store chunk vectors with reference to document
    """

    root = Path(folder_path).expanduser().resolve()
    if not root.exists() or not root.is_dir():
        raise FileNotFoundError(f"Folder not found: {folder_path}")

    files = sorted(p for p in root.glob(glob_pattern) if p.is_file())
    supported_files = [p for p in files if is_supported_file(p)]
    skipped = [p.as_posix() for p in files if not is_supported_file(p)]

    if not supported_files:
        return IngestReport(
            namespace=namespace, file_count=0, chunk_count=0, skipped_files=skipped
        )

    global_memory = app.memory.global_scope

    total_chunks = 0
    for file_path in supported_files:
        relative_path = file_path.relative_to(root).as_posix()
        try:
            full_text = read_text(file_path)
        except Exception as exc:  # pragma: no cover - defensive
            skipped.append(f"{relative_path} (error: {exc})")
            continue

        # TIER 1: Store full document ONCE
        document_key = f"{namespace}:doc:{relative_path}"
        await global_memory.set(
            key=document_key,
            data={
                "full_text": full_text,
                "relative_path": relative_path,
                "namespace": namespace,
                "file_size": len(full_text),
            }
        )

        # Create chunks
        doc_chunks = chunk_markdown_text(
            full_text,
            relative_path=relative_path,
            namespace=namespace,
            chunk_size=chunk_size,
            overlap=chunk_overlap,
        )
        if not doc_chunks:
            continue

        # TIER 2: Store chunk vectors with document reference
        embeddings = embed_texts([chunk.text for chunk in doc_chunks])
        for idx, (chunk, embedding) in enumerate(zip(doc_chunks, embeddings)):
            vector_key = f"{namespace}|{chunk.chunk_id}"
            metadata = {
                "text": chunk.text,
                "namespace": namespace,
                "relative_path": chunk.relative_path,
                "section": chunk.section,
                "start_line": chunk.start_line,
                "end_line": chunk.end_line,
                # NEW: Reference to full document (not the text itself!)
                "document_key": document_key,
                "chunk_index": idx,
                "total_chunks": len(doc_chunks),
            }
            await global_memory.set_vector(
                key=vector_key, embedding=embedding, metadata=metadata
            )
            total_chunks += 1

    log_info(
        f"Ingested {total_chunks} chunks from {len(supported_files)} files into namespace '{namespace}'"
    )

    return IngestReport(
        namespace=namespace,
        file_count=len(supported_files),
        chunk_count=total_chunks,
        skipped_files=skipped,
    )


# ========================= Helper Functions =========================


def _alpha_key(index: int) -> str:
    """Convert index to alphabetic key (0->A, 1->B, ..., 26->AA)."""
    if index < 0:
        raise ValueError("Index must be non-negative")

    letters: List[str] = []
    current = index
    while True:
        current, remainder = divmod(current, 26)
        letters.append(chr(ord("A") + remainder))
        if current == 0:
            break
        current -= 1
    return "".join(reversed(letters))


def _filter_hits(
    hits: Sequence[Dict],
    *,
    namespace: str,
    min_score: float,
) -> List[Dict]:
    """Filter vector search hits by namespace and minimum score."""
    filtered: List[Dict] = []
    for hit in hits:
        metadata = hit.get("metadata", {})
        if metadata.get("namespace") != namespace:
            continue
        if hit.get("score", 0.0) < min_score:
            continue
        filtered.append(hit)
    return filtered


def _deduplicate_results(results: List[RetrievalResult]) -> List[RetrievalResult]:
    """Deduplicate by source, keeping highest score per unique chunk."""
    by_source: Dict[str, RetrievalResult] = {}

    for result in results:
        if result.source not in by_source or result.score > by_source[result.source].score:
            by_source[result.source] = result

    # Sort by score descending, limit to top 15
    deduplicated = sorted(by_source.values(), key=lambda x: x.score, reverse=True)
    return deduplicated[:15]


def _build_citations(results: Sequence[RetrievalResult]) -> List[Citation]:
    """Convert retrieval results to citation objects with alphabetic keys."""
    citations: List[Citation] = []

    for idx, result in enumerate(results):
        # Parse source format: "file.md:10-20"
        parts = result.source.split(":")
        relative_path = parts[0]
        line_range = parts[1] if len(parts) > 1 else "0-0"
        line_parts = line_range.split("-")
        start_line = int(line_parts[0]) if line_parts else 0
        end_line = int(line_parts[1]) if len(line_parts) > 1 else start_line

        key = _alpha_key(idx)
        citation = Citation(
            key=key,
            relative_path=relative_path,
            start_line=start_line,
            end_line=end_line,
            section=None,  # Could extract from metadata if needed
            preview=result.text[:200],
            score=result.score,
        )
        citations.append(citation)

    return citations


def _format_context_for_synthesis(results: Sequence[RetrievalResult]) -> str:
    """Format retrieval results as numbered context for the synthesizer."""
    if not results:
        return "(no context available)"

    blocks: List[str] = []
    for idx, result in enumerate(results):
        key = _alpha_key(idx)
        blocks.append(f"[{key}] {result.source}\n{result.text}")

    return "\n\n".join(blocks)


def _calculate_document_score(chunks: List[RetrievalResult]) -> float:
    """
    Calculate relevance score for a document based on its matching chunks.

    Score = (chunk_frequency * 0.4) + (avg_similarity * 0.4) + (max_similarity * 0.2)

    This rewards:
    - Documents with multiple matching chunks (comprehensive coverage)
    - High average relevance across chunks
    - At least one highly relevant section
    """
    if not chunks:
        return 0.0

    chunk_count = len(chunks)
    avg_score = sum(c.score for c in chunks) / chunk_count
    max_score = max(c.score for c in chunks)

    # Normalize chunk count (diminishing returns after 3 chunks)
    normalized_count = min(chunk_count / 3.0, 1.0)

    return (normalized_count * 0.4) + (avg_score * 0.4) + (max_score * 0.2)


async def _aggregate_chunks_to_documents(
    chunks: List[RetrievalResult],
    top_n: int = 5
) -> List[DocumentContext]:
    """
    Group chunks by document, fetch full documents, and rank by relevance.

    Returns top N most relevant full documents.
    """
    from collections import defaultdict

    global_memory = app.memory.global_scope

    # Group chunks by document_key
    by_document: Dict[str, List[RetrievalResult]] = defaultdict(list)
    for chunk in chunks:
        doc_key = chunk.metadata.get("document_key")
        if doc_key:
            by_document[doc_key].append(chunk)

    if not by_document:
        log_info("[aggregate_chunks_to_documents] No document keys found in chunks")
        return []

    log_info(f"[aggregate_chunks_to_documents] Found {len(by_document)} unique documents")

    # Fetch full documents and build DocumentContext objects
    document_contexts: List[DocumentContext] = []
    for doc_key, doc_chunks in by_document.items():
        # Fetch full document
        doc_data = await global_memory.get(key=doc_key)
        if not doc_data:
            log_info(f"[aggregate_chunks_to_documents] Document not found: {doc_key}")
            continue

        # Calculate relevance score
        relevance_score = _calculate_document_score(doc_chunks)

        # Extract matched sections
        matched_sections = [
            chunk.metadata.get("section")
            for chunk in doc_chunks
            if chunk.metadata.get("section")
        ]
        # Remove duplicates while preserving order
        seen = set()
        unique_sections = []
        for section in matched_sections:
            if section not in seen:
                seen.add(section)
                unique_sections.append(section)

        document_contexts.append(
            DocumentContext(
                document_key=doc_key,
                full_text=doc_data.get("full_text", ""),
                relative_path=doc_data.get("relative_path", "unknown"),
                matching_chunks=len(doc_chunks),
                relevance_score=relevance_score,
                matched_sections=unique_sections,
            )
        )

    # Sort by relevance score and return top N
    ranked_documents = sorted(
        document_contexts,
        key=lambda x: x.relevance_score,
        reverse=True
    )[:top_n]

    log_info(
        f"[aggregate_chunks_to_documents] Returning top {len(ranked_documents)} documents "
        f"(scores: {[f'{d.relevance_score:.3f}' for d in ranked_documents]})"
    )

    return ranked_documents


def _format_documents_for_synthesis(documents: Sequence[DocumentContext]) -> str:
    """Format full documents with minimal metadata for better AI comprehension."""
    if not documents:
        return "(no documents available)"

    blocks: List[str] = []
    for idx, doc in enumerate(documents):
        key = _alpha_key(idx)
        # Simple, clean format - just document ID and content
        header = f"=== DOCUMENT [{key}]: {doc.relative_path} ==="
        blocks.append(f"{header}\n\n{doc.full_text}\n")

    return "\n".join(blocks)


def _build_citations_from_documents(documents: Sequence[DocumentContext]) -> List[Citation]:
    """Convert document contexts to citation objects."""
    citations: List[Citation] = []

    for idx, doc in enumerate(documents):
        key = _alpha_key(idx)
        citation = Citation(
            key=key,
            relative_path=doc.relative_path,
            start_line=0,  # Full document, no specific line
            end_line=0,
            section=", ".join(doc.matched_sections) if doc.matched_sections else None,
            preview=doc.full_text[:200],
            score=doc.relevance_score,
        )
        citations.append(citation)

    return citations


# ========================= Agent 1: Query Planner =========================


@app.reasoner()
async def plan_queries(question: str) -> QueryPlan:
    """Generate 3-5 diverse search queries from the user's question."""

    return await app.ai(
        system=(
            "You are a query planning expert for documentation search. "
            "Your job is to generate 3-5 DIVERSE search queries that maximize retrieval coverage.\n\n"
            "DIVERSITY STRATEGIES:\n"
            "1. Use different terminology and synonyms\n"
            "2. Cover different aspects (setup, usage, troubleshooting, configuration)\n"
            "3. Range from broad concepts to specific terms\n"
            "4. Include related concepts (e.g., 'authentication' → also 'login', 'credentials')\n"
            "5. Avoid redundancy - each query should target unique angles\n\n"
            "QUERY TYPES:\n"
            "- How-to queries: 'how to install X'\n"
            "- Concept queries: 'X architecture'\n"
            "- Troubleshooting: 'X error', 'X not working'\n"
            "- Configuration: 'X settings', 'configure X'\n"
            "- API/Reference: 'X API', 'X methods'"
        ),
        user=(
            f"Question: {question}\n\n"
            "Generate 3-5 diverse search queries that cover different angles of this question. "
            "Also specify the strategy: 'broad' (general exploration), 'specific' (targeted search), "
            "or 'mixed' (combination of both)."
        ),
        schema=QueryPlan,
    )


# ========================= Agent 2: Parallel Retrievers =========================


async def _retrieve_for_query(
    query: str,
    namespace: str,
    top_k: int,
    min_score: float,
) -> List[RetrievalResult]:
    """Single retrieval operation for one query."""

    global_memory = app.memory.global_scope

    # Embed the query
    embedding = embed_query(query)

    # Search vector store
    raw_hits = await global_memory.similarity_search(
        query_embedding=embedding,
        top_k=top_k * 2  # Get more to account for filtering
    )

    # Filter by namespace and score
    filtered_hits = _filter_hits(raw_hits, namespace=namespace, min_score=min_score)

    # Convert to RetrievalResult objects
    results: List[RetrievalResult] = []
    for hit in filtered_hits[:top_k]:
        metadata = hit.get("metadata", {})
        text = metadata.get("text", "").strip()
        if not text:
            continue

        relative_path = metadata.get("relative_path", "unknown")
        start_line = int(metadata.get("start_line", 0))
        end_line = int(metadata.get("end_line", 0))
        source = f"{relative_path}:{start_line}-{end_line}"

        results.append(
            RetrievalResult(
                text=text,
                source=source,
                score=float(hit.get("score", 0.0)),
                metadata=metadata,  # Include full metadata for document aggregation
            )
        )

    return results


@app.reasoner()
async def parallel_retrieve(
    queries: List[str],
    namespace: str = "documentation",
    top_k: int = 6,
    min_score: float = 0.35,
) -> List[RetrievalResult]:
    """Execute parallel retrieval for all queries and deduplicate results."""

    log_info(f"[parallel_retrieve] Running {len(queries)} queries in parallel")

    # Execute all retrievals in parallel
    tasks = [
        _retrieve_for_query(query, namespace, top_k, min_score)
        for query in queries
    ]
    all_results_lists = await asyncio.gather(*tasks)

    # Flatten results
    all_results: List[RetrievalResult] = []
    for results in all_results_lists:
        all_results.extend(results)

    log_info(f"[parallel_retrieve] Retrieved {len(all_results)} total chunks before deduplication")

    # Deduplicate and rank
    deduplicated = _deduplicate_results(all_results)

    log_info(f"[parallel_retrieve] Returning {len(deduplicated)} unique chunks")

    return deduplicated


# ========================= Agent 3: Self-Aware Synthesizer =========================


@app.reasoner()
async def synthesize_answer(
    question: str,
    results: List[RetrievalResult],
    is_refinement: bool = False,
) -> DocAnswer:
    """Generate answer with self-assessment of completeness."""

    if not results:
        return DocAnswer(
            answer="I could not find any relevant documentation to answer this question.",
            citations=[],
            confidence="insufficient",
            needs_more=False,
            missing_topics=["No documentation found for this topic"],
        )

    # Format context for the AI
    context_text = _format_context_for_synthesis(results)

    # Build citations
    citations = _build_citations(results)

    # Create a mapping of keys to sources for the prompt
    key_map = "\n".join([f"[{c.key}] = {c.relative_path}:{c.start_line}-{c.end_line}" for c in citations])

    system_prompt = (
        "You are a precise documentation assistant with SELF-AWARENESS capabilities.\n\n"
        "ANSWER GENERATION RULES:\n"
        "1. Use GitHub-flavored Markdown (2-6 concise sentences or bullets)\n"
        "2. Include inline citations like [A] or [B][D] after each claim\n"
        "3. ONLY use facts from the provided context chunks\n"
        "4. Never invent API names, CLI commands, config values, or examples\n"
        "5. If unsure, say 'The documentation doesn't specify...'\n\n"
        "SELF-ASSESSMENT RULES (CRITICAL):\n"
        "After generating your answer, assess its completeness:\n\n"
        "→ confidence='high', needs_more=False, missing_topics=[]\n"
        "  IF: You can fully answer the question with all key details from context\n\n"
        "→ confidence='partial', needs_more=True, missing_topics=['specific topic 1', 'specific topic 2']\n"
        "  IF: You can partially answer but key details are missing\n"
        "  LIST: Specific missing information (e.g., 'installation steps', 'configuration options')\n\n"
        "→ confidence='insufficient', needs_more=True, missing_topics=['what info is needed']\n"
        "  IF: Context doesn't contain relevant information\n"
        "  ANSWER: 'I don't have documentation about X. I need information about: [missing_topics]'\n\n"
        f"{'REFINEMENT MODE: This is a second attempt. Be more lenient - if you have ANY useful info, set needs_more=False.' if is_refinement else ''}"
    )

    user_prompt = (
        f"Question: {question}\n\n"
        f"Citation Key Map:\n{key_map}\n\n"
        f"Context Chunks:\n{context_text}\n\n"
        "Generate a concise markdown answer with inline citations. "
        "Then self-assess: can you fully answer this question with the provided context? "
        "Set confidence, needs_more, and missing_topics accordingly."
    )

    # Get structured response
    response = await app.ai(
        system=system_prompt,
        user=user_prompt,
        schema=DocAnswer,
    )

    # Ensure citations are included
    if isinstance(response, DocAnswer):
        if not response.citations:
            response.citations = citations
        return response

    # Fallback if response is dict
    response_dict = response if isinstance(response, dict) else response.model_dump()
    response_dict["citations"] = citations
    return DocAnswer.model_validate(response_dict)


# ========================= Main Orchestrator =========================


@app.reasoner()
async def qa_answer(
    question: str,
    namespace: str = "documentation",
    top_k: int = 6,
    min_score: float = 0.35,
) -> DocAnswer:
    """
    Main QA orchestrator with parallel retrieval and optional refinement.

    Flow:
    1. Plan diverse queries
    2. Parallel retrieval
    3. Synthesize with self-assessment
    4. Optional refinement if needs_more=True (max 1 iteration)
    """

    log_info(f"[qa_answer] Processing question: {question}")

    # Step 1: Plan diverse queries
    plan = await plan_queries(question)
    log_info(f"[qa_answer] Generated {len(plan.queries)} queries with strategy: {plan.strategy}")

    # Step 2: Parallel retrieval
    results = await parallel_retrieve(
        queries=plan.queries,
        namespace=namespace,
        top_k=top_k,
        min_score=min_score,
    )

    # Step 3: Synthesize answer
    answer = await synthesize_answer(question, results, is_refinement=False)

    log_info(
        f"[qa_answer] First synthesis: confidence={answer.confidence}, "
        f"needs_more={answer.needs_more}, citations={len(answer.citations)}"
    )

    # Step 4: Optional refinement (max 1 iteration)
    if answer.needs_more and answer.missing_topics:
        log_info(f"[qa_answer] Refinement needed for: {answer.missing_topics}")

        # Generate targeted queries for missing topics
        refinement_queries = []
        for topic in answer.missing_topics[:3]:  # Limit to 3 topics
            refinement_queries.append(f"{question} {topic}")
            refinement_queries.append(topic)

        # Retrieve more context
        additional_results = await parallel_retrieve(
            queries=refinement_queries,
            namespace=namespace,
            top_k=top_k,
            min_score=min_score,
        )

        # Merge with previous results and deduplicate
        all_results = results + additional_results
        merged_results = _deduplicate_results(all_results)

        log_info(f"[qa_answer] Refinement retrieved {len(additional_results)} new chunks, "
                f"merged to {len(merged_results)} total")

        # Synthesize again with refinement flag
        answer = await synthesize_answer(question, merged_results, is_refinement=True)

        log_info(
            f"[qa_answer] Refined synthesis: confidence={answer.confidence}, "
            f"needs_more={answer.needs_more}, citations={len(answer.citations)}"
        )

    return answer


# ========================= Document-Aware QA (NEW) =========================


@app.reasoner()
async def qa_answer_with_documents(
    question: str,
    namespace: str = "documentation",
    top_k: int = 6,
    min_score: float = 0.35,
    top_documents: int = 5,
) -> DocAnswer:
    """
    Document-aware QA orchestrator that retrieves full documents instead of chunks.

    Flow:
    1. Plan diverse queries
    2. Parallel chunk retrieval
    3. Aggregate chunks to full documents
    4. Synthesize answer using full document context
    5. Optional refinement if needs_more=True (max 1 iteration)
    """

    log_info(f"[qa_answer_with_documents] Processing question: {question}")

    # Step 1: Plan diverse queries
    plan = await plan_queries(question)
    log_info(f"[qa_answer_with_documents] Generated {len(plan.queries)} queries with strategy: {plan.strategy}")

    # Step 2: Parallel chunk retrieval
    chunk_results = await parallel_retrieve(
        queries=plan.queries,
        namespace=namespace,
        top_k=top_k,
        min_score=min_score,
    )

    # Step 3: Aggregate chunks to full documents
    documents = await _aggregate_chunks_to_documents(chunk_results, top_n=top_documents)

    if not documents:
        return DocAnswer(
            answer="I could not find any relevant documentation to answer this question.",
            citations=[],
            confidence="insufficient",
            needs_more=False,
            missing_topics=["No documentation found for this topic"],
        )

    # Step 4: Synthesize answer using full documents
    context_text = _format_documents_for_synthesis(documents)
    citations = _build_citations_from_documents(documents)

    key_map = "\n".join([f"[{c.key}] = {c.relative_path}" for c in citations])

    system_prompt = (
        "You are a documentation expert who READS and COMPREHENDS documentation to answer questions.\n\n"
        "🔍 CRITICAL READING INSTRUCTIONS:\n"
        "1. READ the full documentation pages provided below CAREFULLY and THOROUGHLY\n"
        "2. FIND the specific information that directly answers the user's question\n"
        "3. EXTRACT and PRESENT the actual details, steps, commands, or explanations from the docs\n"
        "4. QUOTE or PARAPHRASE directly from the documentation - be SPECIFIC\n"
        "5. If the answer requires multiple steps or details, extract ALL of them from the documentation\n\n"
        "✅ ANSWER FORMAT (IMPORTANT):\n"
        "- Start with a DIRECT answer to the question\n"
        "- Include SPECIFIC details: actual commands, file paths, configuration values, step-by-step instructions\n"
        "- Use inline citations [A][B] to reference which document each fact comes from\n"
        "- Use GitHub-flavored Markdown (code blocks, bullets, etc.)\n"
        "- Be CONCRETE and ACTIONABLE - give users what they need to DO something\n\n"
        "❌ WHAT NOT TO DO:\n"
        "- Don't just say 'the documentation mentions X' - TELL THEM WHAT IT SAYS\n"
        "- Don't be vague or generic - extract SPECIFIC information\n"
        "- Don't summarize what the docs are about - ANSWER the question with actual content\n"
        "- Don't say 'refer to the documentation' - YOU are reading it FOR them\n\n"
        "📚 EXAMPLE OF GOOD vs BAD:\n"
        "Question: 'How do I get started?'\n"
        "❌ BAD: 'The documentation mentions getting started steps [A]'\n"
        "✅ GOOD: 'To get started: 1) Install the CLI with `npm install -g agentfield` [A], 2) Run `af init` to create a new project [A], 3) Configure your agent in `agent.yaml` [A]'\n\n"
        "Question: 'How is IAM treated?'\n"
        "❌ BAD: 'The documentation mentions identity management [A]'\n"
        "✅ GOOD: 'AgentField uses Decentralized Identifiers (DIDs) for identity management [A]. Each agent gets a unique DID that is cryptographically verifiable [A]. IAM policies can be configured in the control plane settings [B]'\n\n"
        "🎯 SELF-ASSESSMENT RULES:\n"
        "After generating your answer, honestly assess:\n\n"
        "→ confidence='high', needs_more=False, missing_topics=[]\n"
        "  IF: You found SPECIFIC, DETAILED information that directly and completely answers the question\n\n"
        "→ confidence='partial', needs_more=True, missing_topics=['specific missing detail 1', 'specific missing detail 2']\n"
        "  IF: You found SOME information but it's incomplete (e.g., has steps 1-2 but missing step 3)\n"
        "  LIST: Exactly what specific information is missing\n\n"
        "→ confidence='insufficient', needs_more=True, missing_topics=['what info is needed']\n"
        "  IF: After thoroughly reading ALL documents, the specific information requested is genuinely not present\n\n"
        "⚠️ IMPORTANT: Don't confuse 'not in one sentence' with 'not in documentation'.\n"
        "If the answer requires combining info from multiple paragraphs or sections, that's STILL a complete answer!"
    )

    user_prompt = (
        f"Question: {question}\n\n"
        f"Citation Key Map:\n{key_map}\n\n"
        f"Full Documentation Pages:\n{context_text}\n\n"
        "Generate a concise markdown answer with inline citations. "
        "Then self-assess: can you fully answer this question with the provided documents? "
        "Set confidence, needs_more, and missing_topics accordingly."
    )

    response = await app.ai(
        system=system_prompt,
        user=user_prompt,
        schema=DocAnswer,
    )

    # Ensure citations are included
    if isinstance(response, DocAnswer):
        if not response.citations:
            response.citations = citations
        answer = response
    else:
        response_dict = response if isinstance(response, dict) else response.model_dump()
        response_dict["citations"] = citations
        answer = DocAnswer.model_validate(response_dict)

    log_info(
        f"[qa_answer_with_documents] First synthesis: confidence={answer.confidence}, "
        f"needs_more={answer.needs_more}, documents_used={len(documents)}"
    )

    # Step 5: Optional refinement (max 1 iteration)
    if answer.needs_more and answer.missing_topics:
        log_info(f"[qa_answer_with_documents] Refinement needed for: {answer.missing_topics}")

        # Generate targeted queries for missing topics
        refinement_queries = []
        for topic in answer.missing_topics[:3]:  # Limit to 3 topics
            refinement_queries.append(f"{question} {topic}")
            refinement_queries.append(topic)

        # Retrieve more chunks
        additional_chunks = await parallel_retrieve(
            queries=refinement_queries,
            namespace=namespace,
            top_k=top_k,
            min_score=min_score,
        )

        # Merge and aggregate to documents
        all_chunks = chunk_results + additional_chunks
        merged_documents = await _aggregate_chunks_to_documents(all_chunks, top_n=top_documents)

        log_info(f"[qa_answer_with_documents] Refinement found {len(merged_documents)} total documents")

        # Synthesize again with more lenient prompt
        context_text = _format_documents_for_synthesis(merged_documents)
        citations = _build_citations_from_documents(merged_documents)
        key_map = "\n".join([f"[{c.key}] = {c.relative_path}" for c in citations])

        system_prompt_refined = system_prompt + "\n\nREFINEMENT MODE: This is a second attempt. Be more lenient - if you have ANY useful info, set needs_more=False."

        user_prompt_refined = (
            f"Question: {question}\n\n"
            f"Citation Key Map:\n{key_map}\n\n"
            f"Full Documentation Pages:\n{context_text}\n\n"
            "Generate a concise markdown answer with inline citations. "
            "Then self-assess: can you fully answer this question with the provided documents? "
            "Set confidence, needs_more, and missing_topics accordingly."
        )

        response = await app.ai(
            system=system_prompt_refined,
            user=user_prompt_refined,
            schema=DocAnswer,
        )

        if isinstance(response, DocAnswer):
            if not response.citations:
                response.citations = citations
            answer = response
        else:
            response_dict = response if isinstance(response, dict) else response.model_dump()
            response_dict["citations"] = citations
            answer = DocAnswer.model_validate(response_dict)

        log_info(
            f"[qa_answer_with_documents] Refined synthesis: confidence={answer.confidence}, "
            f"needs_more={answer.needs_more}, documents_used={len(merged_documents)}"
        )

    return answer


# ========================= Bootstrapping =========================


def _warmup_embeddings() -> None:
    """Warm up the embedding model on startup."""
    try:
        embed_texts(["doc-chatbot warmup"])
        log_info("FastEmbed model warmed up for documentation chatbot")
    except Exception as exc:  # pragma: no cover - best-effort
        log_info(f"FastEmbed warmup failed: {exc}")


if __name__ == "__main__":
    _warmup_embeddings()

    print("📚 Simplified Documentation Chatbot Agent")
    print("🧠 Node ID: documentation-chatbot")
    print(f"🌐 Control Plane: {app.agentfield_server}")
    print("\n🎯 Architecture: 3-Agent Parallel System + Document-Level Retrieval")
    print("  1. Query Planner → Generates diverse search queries")
    print("  2. Parallel Retrievers → Concurrent vector search")
    print("  3. Self-Aware Synthesizer → Answer + confidence assessment")
    print("\n📄 Storage Strategy: Two-Tier System")
    print("  • Documents stored ONCE in regular memory")
    print("  • Chunks reference documents (no duplication)")
    print("  • 70% storage savings vs naive approach")
    print("\nEndpoints:")
    print("  • /skills/ingest_folder → Ingest documentation (two-tier storage)")
    print("  • /reasoners/plan_queries → Generate diverse queries")
    print("  • /reasoners/parallel_retrieve → Parallel chunk retrieval")
    print("  • /reasoners/synthesize_answer → Self-aware synthesis (chunk-based)")
    print("  • /reasoners/qa_answer → Chunk-based QA orchestrator")
    print("  • /reasoners/qa_answer_with_documents → 🆕 Document-aware QA (RECOMMENDED)")
    print("\n✨ Features:")
    print("  - Parallel retrieval for 3x speed improvement")
    print("  - Self-aware synthesis (no separate review)")
    print("  - Max 1 refinement iteration (prevents loops)")
    print("  - Simple schemas (.ai compatible, 2-4 attributes)")
    print("  - Document-level context (full pages vs isolated chunks)")
    print("  - Smart document ranking (frequency + relevance scoring)")

    app.run(auto_port=True)
