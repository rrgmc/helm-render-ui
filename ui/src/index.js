import React from "react";
import ReactDOM from "react-dom";
import HelmTemplatePreview from "./components/helm";
import config from "./config";

import "./styles/base.css";
import "./styles/navbar.css";
import "./styles/editor.css";

ReactDOM.render(
  <React.StrictMode>
    <HelmTemplatePreview apiURL={config.api} />
  </React.StrictMode>,
  document.getElementById("root")
);
