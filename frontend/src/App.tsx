import React from "react";
import SortPage from "./pages/SortPage";
import "./App.css";

function App() {
  return (
    <div
      style={{
        backgroundColor: "#1f1f1f",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100vh",
        width: "100vw",
      }}
    >
      <SortPage />
    </div>
  );
}

export default App;
