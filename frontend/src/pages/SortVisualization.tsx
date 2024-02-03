import React from "react";
import Box from "@mui/material/Box";
import Paper from "@mui/material/Paper";
import { SortData } from "./SortPage";

interface SortVisualizationProps {
  sortData: SortData;
}

const SortVisualization = ({ sortData }: SortVisualizationProps) => {
  const maxValue = Math.max(...sortData.data);

  return (
    <Box
      sx={{
        flexDirection: "row",
        display: "flex",
        alignItems: "flex-end",
        height: 300,
        width: 800,
        padding: 2,
        overflowX: "auto",
      }}
    >
      {sortData.data.map((item, index) => (
        <Paper
          key={index}
          sx={{
            width: 800 / sortData.data.length,
            height: `${(item / maxValue) * 100}%`,
            display: "flex",
            justifyContent: "center",
            alignItems: "flex-end",
            borderBottomLeftRadius: 0,
            borderBottomRightRadius: 0,
            backgroundColor:
              sortData.lastConsidered === index ? "#d32f2f" : "#1976d2",
          }}
        ></Paper>
      ))}
    </Box>
  );
};

export default SortVisualization;
