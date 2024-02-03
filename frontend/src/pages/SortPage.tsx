import React, { useCallback, useEffect, useState } from "react";
import {
  Box,
  Button,
  MenuItem,
  Select,
  Slider,
  TextField,
  Typography,
} from "@mui/material";
import SortVisualization from "./SortVisualization";
import { debounce } from "lodash";

export interface SortData {
  data: number[];
  lastConsidered: number | null;
}

const SortPage = () => {
  const [ws, setWs] = useState<WebSocket>();
  const [selectedAlgorithm, setSelectedAlgorithm] = useState<string>("bubble");
  const [randomCount, setRandomCount] = useState<number>(20);
  const [stepDelay, setStepDelay] = useState<number>(100);
  const [sortData, setSortData] = useState<SortData>({
    data: [],
    lastConsidered: null,
  });

  useEffect(() => {
    randomizeInput();
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
      console.log("connected");
    };

    socket.onmessage = (e) => {
      const data = JSON.parse(e.data);
      if (data.type === "progress") {
        setSortData({
          data: data.data,
          lastConsidered: data.lastConsidered,
        });
      }

      if (data.type === "stepDelayConfirmation") {
        console.log("stepDelayConfirmation", data);
      }
    };

    socket.onclose = () => {
      console.log("disconnected");
    };

    setWs(socket);
    return () => {
      socket.close();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const randomizeInput = () => {
    const randomArray = Array.from(
      { length: randomCount },
      () => Math.floor(Math.random() * randomCount) + 1
    );
    setSortData({
      data: randomArray,
      lastConsidered: null,
    });
  };

  const sendStepDelayUpdate = useCallback(
    (value: number) => {
      if (ws) {
        console.log("sending stepDelay", value);
        ws.send(
          JSON.stringify({
            type: "stepDelay",
            data: value,
          })
        );
      }
    },
    [ws]
  );

  // eslint-disable-next-line react-hooks/exhaustive-deps
  const debouncedSendStepDelayUpdate = useCallback(
    debounce(sendStepDelayUpdate, 200),
    [sendStepDelayUpdate]
  );

  const updateStepDelay = (value: number) => {
    setStepDelay(value);
    debouncedSendStepDelayUpdate(value);
  };

  const sortArray = () => {
    if (ws) {
      ws.send(
        JSON.stringify({
          type: "sort",
          algorithm: selectedAlgorithm,
          data: sortData?.data,
        })
      );
    }
  };

  return (
    <Box
      sx={{
        flexDirection: "column",
        display: "flex",
        alignItems: "center",
      }}
    >
      <Box
        sx={{
          flexDirection: "row",
          display: "flex",
          alignItems: "center",
        }}
      >
        <Typography variant="h5">Sort</Typography>
        <Select
          sx={{
            backgroundColor: "#fff",
            marginLeft: 5,
          }}
          label="Algorithm"
          value={selectedAlgorithm}
          onChange={(e) => setSelectedAlgorithm(e.target.value)}
        >
          <MenuItem value="bubble">Bubble</MenuItem>
          <MenuItem value="merge">Merge</MenuItem>
          <MenuItem value="quick">Quick</MenuItem>
        </Select>
        {/* <TextField
          sx={{ marginLeft: 5, backgroundColor: "#fff" }}
          variant="outlined"
          placeholder="1,2,3,4,5"
          value={dataInput}
          onChange={(e) => setDataInput(e.target.value)}
        /> */}
        <Button sx={{ marginLeft: 5 }} variant="contained" onClick={sortArray}>
          Sort
        </Button>
        <Button
          sx={{ marginLeft: 2 }}
          variant="contained"
          onClick={randomizeInput}
        >
          randomize
        </Button>
        <TextField
          sx={{ marginLeft: 5, backgroundColor: "#fff" }}
          variant="outlined"
          placeholder="20"
          value={randomCount}
          onChange={(e) => setRandomCount(parseInt(e.target.value) || 0)}
        />
      </Box>
      <SortVisualization sortData={sortData} />
      <Slider
        sx={{ width: 400, marginTop: 2 }}
        value={stepDelay}
        onChange={(_, value) => updateStepDelay(value as number)}
        min={0}
        step={0.1}
        max={100}
        valueLabelDisplay="auto"
        valueLabelFormat={(value) => `${value}`}
      />
    </Box>
  );
};

export default SortPage;
