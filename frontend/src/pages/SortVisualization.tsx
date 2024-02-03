import React, { useEffect, useRef } from "react";
import Box from "@mui/material/Box";
import { SortData } from "./SortPage";

interface SortVisualizationProps {
  sortData: SortData;
}

const SortVisualization = ({ sortData }: SortVisualizationProps) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) {
      return;
    }
    const context = canvas.getContext("2d");
    if (!context) {
      return;
    }
    const width = canvas.width;
    const height = canvas.height;

    context.clearRect(0, 0, width, height);

    const barWidth = width / sortData.data.length;
    const maxValue = Math.max(...sortData.data);

    sortData.data.forEach((item, index) => {
      const barHeight = (item / maxValue) * height;
      context.fillStyle =
        sortData.lastConsidered === index ? "#d32f2f" : "#1976d2";
      context.fillRect(
        index * barWidth,
        height - barHeight,
        barWidth,
        barHeight
      );
    });
  }, [sortData]);

  return (
    <Box
      sx={{
        padding: 2,
        overflowX: "auto",
      }}
    >
      <canvas ref={canvasRef} width={800} height={300} />
    </Box>
  );
};

export default SortVisualization;
