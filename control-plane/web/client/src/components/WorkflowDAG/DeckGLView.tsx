import { useEffect, useMemo, useState } from "react";
import DeckGL from "@deck.gl/react";
import {
  COORDINATE_SYSTEM,
  OrthographicController,
  OrthographicView,
  type OrthographicViewState,
} from "@deck.gl/core";
import { ScatterplotLayer, PathLayer } from "@deck.gl/layers";
import type { WorkflowDAGLightweightNode } from "../../types/workflows";

export type WorkflowDAGNode = WorkflowDAGLightweightNode & {
  workflow_id?: string;
};

interface DeckNode {
  id: string;
  position: [number, number, number];
  depth: number;
  radius: number;
  fillColor: [number, number, number, number];
  borderColor: [number, number, number, number];
  glowColor: [number, number, number, number];
  original: WorkflowDAGNode;
}

interface DeckEdge {
  id: string;
  path: [number, number, number][];
  color: [number, number, number, number];
  width: number;
}

export interface AgentPaletteEntry {
  agentId: string;
  label: string;
  color: string;
  background: string;
  text: string;
}

interface WorkflowDeckGLViewProps {
  nodes: DeckNode[];
  edges: DeckEdge[];
  onNodeClick?: (node: WorkflowDAGNode) => void;
  onNodeHover?: (node: WorkflowDAGNode | null) => void;
}

export interface DeckGraphData {
  nodes: DeckNode[];
  edges: DeckEdge[];
  agentPalette: AgentPaletteEntry[];
}

const initialViewState: OrthographicViewState = {
  target: [0, 0, 0],
  zoom: 0,
  maxZoom: 8,
  minZoom: -6,
};

function hashString(input: string): number {
  let hash = 0;
  for (let i = 0; i < input.length; i++) {
    hash = (hash << 5) - hash + input.charCodeAt(i);
    hash |= 0;
  }
  return Math.abs(hash);
}

function hslToRgb(h: number, s: number, l: number): [number, number, number] {
  const sat = s / 100;
  const light = l / 100;

  if (sat === 0) {
    const val = Math.round(light * 255);
    return [val, val, val];
  }

  const hue2rgb = (p: number, q: number, t: number) => {
    if (t < 0) t += 1;
    if (t > 1) t -= 1;
    if (t < 1 / 6) return p + (q - p) * 6 * t;
    if (t < 1 / 2) return q;
    if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
    return p;
  };

  const q = light < 0.5 ? light * (1 + sat) : light + sat - light * sat;
  const p = 2 * light - q;
  const hk = h / 360;

  const r = Math.round(hue2rgb(p, q, hk + 1 / 3) * 255);
  const g = Math.round(hue2rgb(p, q, hk) * 255);
  const b = Math.round(hue2rgb(p, q, hk - 1 / 3) * 255);

  return [r, g, b];
}

function getAgentColor(agentId: string, index: number): {
  rgb: [number, number, number];
  css: string;
} {
  const hash = hashString(agentId || `agent-${index}`);
  const golden = 0.6180339887498948;
  const hue = (hash * golden * 360) % 360;
  const saturation = 65 + (hash % 15); // 65-80%
  const lightness = 50 + (hash % 12); // 50-62%
  const rgb = hslToRgb(hue, saturation, lightness);
  return { rgb, css: `rgb(${rgb.join(",")})` };
}

function mixColor(
  color: [number, number, number],
  target: [number, number, number],
  ratio: number
): [number, number, number] {
  return [
    Math.round(color[0] * ratio + target[0] * (1 - ratio)),
    Math.round(color[1] * ratio + target[1] * (1 - ratio)),
    Math.round(color[2] * ratio + target[2] * (1 - ratio)),
  ];
}

export const WorkflowDeckGLView = ({
  nodes,
  edges,
  onNodeClick,
  onNodeHover,
}: WorkflowDeckGLViewProps) => {
  const [viewState, setViewState] =
    useState<OrthographicViewState>(initialViewState);

  useEffect(() => {
    if (!nodes.length) return;

    const xs = nodes.map((node) => node.position[0]);
    const ys = nodes.map((node) => node.position[1]);
    const minX = Math.min(...xs);
    const maxX = Math.max(...xs);
    const minY = Math.min(...ys);
    const maxY = Math.max(...ys);

    const padding = 100;
    const width = maxX - minX || 1;
    const height = maxY - minY || 1;

    setViewState((prev) => ({
      ...prev,
      target: [minX + width / 2, minY + height / 2, 0],
      zoom: Math.log2(Math.min(1200 / (width + padding), 800 / (height + padding))),
    }));
  }, [nodes]);

  const inverseZoom = useMemo(() => {
    const zoomRaw = viewState.zoom;
    const zoomValue = Array.isArray(zoomRaw)
      ? zoomRaw[0] ?? 0
      : zoomRaw ?? 0;
    const factor = Math.pow(2, zoomValue);
    return factor === 0 ? 1 : 1 / factor;
  }, [viewState.zoom]);

  const layers = useMemo(() => {
    const glowLayer = new ScatterplotLayer<DeckNode>({
      id: "workflow-node-glow",
      data: nodes,
      pickable: false,
      radiusScale: 1,
      radiusMinPixels: 8,
      radiusMaxPixels: 100,
      getPosition: (node) => node.position,
      getRadius: (node) => node.radius * 2.2 * inverseZoom,
      getFillColor: (node) => node.glowColor,
      coordinateSystem: COORDINATE_SYSTEM.CARTESIAN,
      updateTriggers: {
        getRadius: inverseZoom,
      },
    });

    const nodeLayer = new ScatterplotLayer<DeckNode>({
      id: "workflow-nodes",
      data: nodes,
      pickable: true,
      radiusScale: 1,
      radiusMinPixels: 4,
      radiusMaxPixels: 80,
      getPosition: (node) => node.position,
      getRadius: (node) => node.radius * inverseZoom,
      getFillColor: (node) => node.fillColor,
      getLineColor: (node) => node.borderColor,
      getLineWidth: () => 1.5,
      lineWidthMinPixels: 1.5,
      stroked: true,
      autoHighlight: true,
      highlightColor: [255, 255, 255, 200],
      coordinateSystem: COORDINATE_SYSTEM.CARTESIAN,
      updateTriggers: {
        getRadius: inverseZoom,
      },
      onClick: (info) => {
        const deckNode = info.object;
        if (deckNode?.original) {
          onNodeClick?.(deckNode.original);
        }
      },
      onHover: (info) => {
        const deckNode = info.object;
        if (deckNode?.original) {
          onNodeHover?.(deckNode.original);
        } else {
          onNodeHover?.(null);
        }
      },
    });

    const edgeLayer = new PathLayer<DeckEdge>({
      id: "workflow-edges",
      data: edges,
      getPath: (edge) => edge.path,
      getColor: (edge) => edge.color,
      getWidth: (edge) => edge.width,
      widthMinPixels: 1,
      widthMaxPixels: 4,
      widthUnits: "pixels",
      rounded: true,
      miterLimit: 2,
      coordinateSystem: COORDINATE_SYSTEM.CARTESIAN,
    });

    return [edgeLayer, glowLayer, nodeLayer];
  }, [nodes, edges, onNodeClick, onNodeHover, inverseZoom]);

  return (
    <DeckGL
      views={new OrthographicView({})}
      controller={{ type: OrthographicController, inertia: true }}
      viewState={viewState}
      onViewStateChange={({ viewState: next }) =>
        setViewState(next as OrthographicViewState)
      }
      layers={layers}
      style={{ width: "100%", height: "100%" }}
    />
  );
};

const BACKGROUND_RGB: [number, number, number] = [11, 18, 32];

/**
 * Create a smooth cubic Bezier curve between two points
 * Optimized: Only 8 segments instead of 32 for better performance
 */
function createCubicBezier(
  source: [number, number, number],
  target: [number, number, number],
  curvature: number = 0.5
): [number, number, number][] {
  const dx = target[0] - source[0];

  // Control points for smooth S-curve
  const control1: [number, number, number] = [
    source[0] + dx * curvature,
    source[1],
    source[2],
  ];

  const control2: [number, number, number] = [
    target[0] - dx * curvature,
    target[1],
    target[2],
  ];

  // Sample the curve with fewer points for performance
  const points: [number, number, number][] = [];
  const segments = 8; // Reduced from 32 for better performance

  for (let i = 0; i <= segments; i++) {
    const t = i / segments;
    const mt = 1 - t;
    const x =
      mt * mt * mt * source[0] +
      3 * mt * mt * t * control1[0] +
      3 * mt * t * t * control2[0] +
      t * t * t * target[0];
    const y =
      mt * mt * mt * source[1] +
      3 * mt * mt * t * control1[1] +
      3 * mt * t * t * control2[1] +
      t * t * t * target[1];
    const z =
      mt * mt * mt * source[2] +
      3 * mt * mt * t * control1[2] +
      3 * mt * t * t * control2[2] +
      t * t * t * target[2];
    points.push([x, y, z]);
  }

  return points;
}

/**
 * Build a hierarchical DAG layout optimized for scalability
 * Uses layer-based positioning for clear directional flow
 */
export function buildDeckGraph(
  timeline: WorkflowDAGNode[],
  horizontalSpacing: number = 200,
  verticalSpacing: number = 80
): DeckGraphData {
  if (!timeline.length) {
    return { nodes: [], edges: [], agentPalette: [] };
  }

  const nodeById = new Map<string, WorkflowDAGNode>();
  const childrenByParent = new Map<string, WorkflowDAGNode[]>();
  const parentsByChild = new Map<string, string[]>();
  const agentColors = new Map<
    string,
    { rgb: [number, number, number]; css: string }
  >();

  // Build graph structure
  timeline.forEach((node, index) => {
    nodeById.set(node.execution_id, node);

    if (node.parent_execution_id) {
      // Track children
      if (!childrenByParent.has(node.parent_execution_id)) {
        childrenByParent.set(node.parent_execution_id, []);
      }
      childrenByParent.get(node.parent_execution_id)!.push(node);

      // Track parents
      if (!parentsByChild.has(node.execution_id)) {
        parentsByChild.set(node.execution_id, []);
      }
      parentsByChild.get(node.execution_id)!.push(node.parent_execution_id);
    }

    // Generate agent colors
    const agentId = node.agent_node_id || `agent-${index}`;
    if (!agentColors.has(agentId)) {
      agentColors.set(agentId, getAgentColor(agentId, agentColors.size));
    }
  });

  if (process.env.NODE_ENV !== "production") {
    console.debug(
      "[DeckGL] Processing",
      timeline.length,
      "nodes with",
      childrenByParent.size,
      "parent nodes"
    );
  }

  // Find root nodes (nodes with no parents or parents not in the graph)
  const roots = timeline.filter(
    (node) =>
      !node.parent_execution_id ||
      !nodeById.has(node.parent_execution_id)
  );

  if (roots.length === 0 && timeline.length > 0) {
    // Fallback: use node with smallest depth
    const fallbackRoot = timeline.reduce((best, node) => {
      const depth = node.workflow_depth ?? Infinity;
      const bestDepth = best.workflow_depth ?? Infinity;
      return depth < bestDepth ? node : best;
    });
    roots.push(fallbackRoot);
  }

  if (process.env.NODE_ENV !== "production") {
    console.debug("[DeckGL] Found", roots.length, "root nodes");
  }

  // Assign nodes to layers using topological sort (BFS)
  const layers: WorkflowDAGNode[][] = [];
  const nodeToLayer = new Map<string, number>();
  const visited = new Set<string>();
  const queue: { node: WorkflowDAGNode; layer: number }[] = [];

  // Initialize with roots
  roots.forEach((root) => {
    queue.push({ node: root, layer: 0 });
  });

  while (queue.length > 0) {
    const { node, layer } = queue.shift()!;

    if (visited.has(node.execution_id)) {
      // If already visited, potentially update layer if this path is longer
      const currentLayer = nodeToLayer.get(node.execution_id)!;
      if (layer > currentLayer) {
        // Remove from old layer
        const oldLayerNodes = layers[currentLayer];
        const idx = oldLayerNodes.findIndex(n => n.execution_id === node.execution_id);
        if (idx >= 0) oldLayerNodes.splice(idx, 1);

        // Add to new layer
        nodeToLayer.set(node.execution_id, layer);
        if (!layers[layer]) layers[layer] = [];
        layers[layer].push(node);
      }
      continue;
    }

    visited.add(node.execution_id);
    nodeToLayer.set(node.execution_id, layer);

    if (!layers[layer]) {
      layers[layer] = [];
    }
    layers[layer].push(node);

    // Add children to next layer
    const children = childrenByParent.get(node.execution_id) ?? [];
    children.forEach((child) => {
      queue.push({ node: child, layer: layer + 1 });
    });
  }

  if (process.env.NODE_ENV !== "production") {
    console.debug(
      "[DeckGL] Created",
      layers.length,
      "layers, max layer size:",
      Math.max(...layers.map(l => l.length))
    );
  }

  // Position nodes in each layer
  const layoutInfo = new Map<
    string,
    {
      position: [number, number, number];
      layer: number;
      color: [number, number, number];
      agentId: string;
      radius: number;
    }
  >();

  layers.forEach((layerNodes, layerIndex) => {
    const layerHeight = (layerNodes.length - 1) * verticalSpacing;
    const startY = -layerHeight / 2;

    layerNodes.forEach((node, indexInLayer) => {
      const x = layerIndex * horizontalSpacing;
      const y = startY + indexInLayer * verticalSpacing;
      const z = 0; // Keep flat for now

      const agentId = node.agent_node_id || node.reasoner_id || "agent";
      const colorInfo =
        agentColors.get(agentId) ??
        getAgentColor(agentId, agentColors.size + 1);
      const baseColor = colorInfo.rgb;

      // Node size based on depth (earlier nodes slightly larger)
      const baseRadius = Math.max(6, 12 - layerIndex * 0.3);

      layoutInfo.set(node.execution_id, {
        position: [x, y, z],
        layer: layerIndex,
        color: baseColor,
        agentId,
        radius: baseRadius,
      });
    });
  });

  // Create deck nodes
  const deckNodes: DeckNode[] = [];
  layoutInfo.forEach((info, nodeId) => {
    const fill = [...info.color, 240] as [number, number, number, number];
    const border = [...mixColor(info.color, BACKGROUND_RGB, 0.4), 255] as [
      number,
      number,
      number,
      number,
    ];
    const glow = [...mixColor(info.color, [255, 255, 255], 0.25), 90] as [
      number,
      number,
      number,
      number,
    ];

    deckNodes.push({
      id: nodeId,
      position: info.position,
      depth: info.layer,
      radius: info.radius,
      fillColor: fill,
      borderColor: border,
      glowColor: glow,
      original: nodeById.get(nodeId)!,
    });
  });

  // Create edges with smooth curves
  const deckEdges: DeckEdge[] = [];
  timeline.forEach((node) => {
    if (!node.parent_execution_id) {
      return;
    }

    const parentInfo = layoutInfo.get(node.parent_execution_id);
    const childInfo = layoutInfo.get(node.execution_id);
    if (!parentInfo || !childInfo) {
      return;
    }

    const source = parentInfo.position;
    const target = childInfo.position;

    // Calculate curvature based on distance
    const dx = Math.abs(target[0] - source[0]);
    const curvature = Math.min(0.6, 0.3 + dx / 1000);

    const path = createCubicBezier(source, target, curvature);

    const baseColor = parentInfo.color;
    const edgeColor = [
      ...mixColor(baseColor, BACKGROUND_RGB, 0.5),
      140,
    ] as [number, number, number, number];

    // Edge width decreases with depth
    const width = Math.max(1.2, 2.5 - childInfo.layer * 0.15);

    deckEdges.push({
      id: `${node.parent_execution_id}-${node.execution_id}`,
      path,
      color: edgeColor,
      width,
    });
  });

  // Build agent palette
  const agentPalette: AgentPaletteEntry[] = [];
  agentColors.forEach((value, key) => {
    const background = mixColor(value.rgb, BACKGROUND_RGB, 0.85);
    const luminance =
      value.rgb[0] * 0.2126 + value.rgb[1] * 0.7152 + value.rgb[2] * 0.0722;
    const textColor = luminance > 140 ? "#0f172a" : "#f8fafc";
    agentPalette.push({
      agentId: key,
      label: key,
      color: value.css,
      background: `rgb(${background.join(",")})`,
      text: textColor,
    });
  });
  agentPalette.sort((a, b) => a.label.localeCompare(b.label));

  if (process.env.NODE_ENV !== "production") {
    console.debug(
      "[DeckGL] Built",
      deckNodes.length,
      "nodes and",
      deckEdges.length,
      "edges"
    );
  }

  return { nodes: deckNodes, edges: deckEdges, agentPalette };
}
