import React from "react";
import Editor from "react-simple-code-editor";
import Preview from "./preview";
import debounce from "lodash.debounce";
import { Tab, Tabs, TabList, TabPanel } from "react-tabs";
import { highlight, languages } from "prismjs/components/prism-core";

import "react-tabs/style/react-tabs.css";
import "prismjs/themes/prism.css";

import "prismjs/components/prism-clike";
import "prismjs/components/prism-yaml";

import { ReactComponent as Logo } from "../static/goggles.svg";

type Props = {
  apiURL: string,
};

export default class HelmTemplatePreview extends React.Component<Props> {
  constructor(props) {
    super(props);
    this.state = {
      rawTemplate: "",
      rawValues: "",
      rawHelpers: "",
      rawChart: defaults.rawChart,
      rawRelease: defaults.rawRelease,
      rawCapabilities: defaults.rawCapabilities,
      renderedTemplate: "",
      renderError: "",
    };
  }

  componentDidMount() {
    this.updateHelmRender();
  }

  updateRawTemplate(rawTemplate) {
    this.setState(
      {
        rawTemplate: rawTemplate,
      },
      this.updateHelmRenderDebounce
    );
  }

  updateRawValues(rawValues) {
    this.setState(
      {
        rawValues: rawValues,
      },
      this.updateHelmRenderDebounce
    );
  }

  updateRawChart(rawChart) {
    this.setState(
      {
        rawChart: rawChart,
      },
      this.updateHelmRenderDebounce
    );
  }

  updateRawRelease(rawRelease) {
    this.setState(
      {
        rawRelease: rawRelease,
      },
      this.updateHelmRenderDebounce
    );
  }

  updateRawCapabilities(rawCapabilities) {
    this.setState(
      {
        rawCapabilities: rawCapabilities,
      },
      this.updateHelmRenderDebounce
    );
  }

  updateRawHelpers(rawHelpers) {
    this.setState(
      {
        rawHelpers: rawHelpers,
      },
      this.updateHelmRenderDebounce
    );
  }

  updateHelmRenderDebounce = debounce(this.updateHelmRender, 300);

  updateHelmRender() {
    const handleResponse = (res) => {
      if (!res.ok) {
        throw res;
      }
      return res;
    };

    const renderTemplate = (res) => {
      res
        .json()
        .then((data) =>
          this.setState({ renderedTemplate: data.preview, renderError: "" })
        );
    };

    const renderError = (error) => {
      error
        .text()
        .then((errorMessage) => this.setState({ renderError: errorMessage }));
    };

    fetch(`${this.props.apiURL}/data`, {
      method: "GET",
    })
      .then(handleResponse)
      .then(renderTemplate)
      .catch(renderError);
  }

  render() {
    const style = {
      whiteSpace: "pre",
      fontFamily: '"Fira code", "Fira Mono", monospace',
      fontSize: 12,
      backgroundColor: "#f5f5f5",
      width: "100%",
    };

    const padding = 12;

    const highlighter = (code) =>
      highlight(code, languages.yaml)
        .split("\n")
        .map(
          (line, idx) =>
            `<span class="editor__line__number">${idx + 1}</span>${line}`
        )
        .join("\n");

    return (
      <div className="app">
        <div className="navbar">
          <Logo title="Helm Template Preview" className="navbar__logo" />
          <h1 className="navbar__title">Helm Template Preview</h1>
          <h3 className="navbar__about">
            <a href="https://zainp.com">About</a>
          </h3>
        </div>
        <div className="container">
          <div className="input">
            <div className="input__values">
              <Tabs>
                <TabList>
                  <Tab>Values</Tab>
                  <Tab>Chart</Tab>
                  <Tab>Release</Tab>
                  <Tab>Capabilities</Tab>
                  <Tab>_helpers.tpl</Tab>
                </TabList>
                <TabPanel>
                  <Editor
                    value={this.state.rawValues}
                    highlight={highlighter}
                    padding={padding}
                    style={style}
                    onValueChange={(code) => this.updateRawValues(code)}
                    placeholder="Insert the contents of your values files"
                    className="input__values__editor editor"
                  />
                </TabPanel>
                <TabPanel>
                  <Editor
                    value={this.state.rawChart}
                    highlight={highlighter}
                    padding={padding}
                    style={style}
                    onValueChange={(code) => this.updateRawChart(code)}
                    placeholder="Insert the contents of your Chart.yaml file (e.g name)"
                    className="input__chart__editor editor"
                  />
                </TabPanel>
                <TabPanel>
                  <Editor
                    value={this.state.rawRelease}
                    highlight={highlighter}
                    padding={padding}
                    style={style}
                    onValueChange={(code) => this.updateRawRelease(code)}
                    placeholder="Insert release values (e.g Name, Namespace, etc...)"
                    className="input__release__editor editor"
                  />
                </TabPanel>
                <TabPanel>
                  <Editor
                    value={this.state.rawCapabilities}
                    highlight={highlighter}
                    padding={padding}
                    style={style}
                    onValueChange={(code) => this.updateRawCapabilities(code)}
                    placeholder="Insert capabilities values (e.g HelmVersion, KubeVersion, etc...)"
                    className="input__capabilities__editor editor"
                  />
                </TabPanel>
                <TabPanel>
                  <Editor
                    value={this.state.rawHelpers}
                    highlight={highlighter}
                    padding={padding}
                    style={style}
                    onValueChange={(code) => this.updateRawHelpers(code)}
                    placeholder="Insert _helpers.tpl file contents"
                    className="input__helpers_tpl__editor editor"
                  />
                </TabPanel>
              </Tabs>
            </div>
          </div>
          <div className="preview">
            <Preview
              value={this.state.renderedTemplate}
              highlight={highlighter}
              padding={padding}
              style={style}
              className="preview__highlighted"
            />
          </div>
        </div>
        <div
          className="render-error"
          style={this.state.renderError === "" ? { display: "none" } : {}}
        >
          {this.state.renderError}
        </div>
      </div>
    );
  }
}

const defaults = {
  rawChart: `apiVersion: v2
name: chart-name
version: 0.1.0`,
  rawRelease: `Name: release-name
Namespace: namespace
IsUpgrade: false
IsInstall: true
Revision: 1
Service: Helm`,
  rawCapabilities: `APIVersions:
  - networking.k8s.io/v1
  - authentication.k8s.io/v1
KubeVersion:
  Version: "1.22"
  Major: "1"
  Minor: "22"
  GitVersion: "1.22.0"
HelmVersion:
  Version: 3.6.3
  GitCommit: d506314abfb5d21419df8c7e7e68012379db2354
  GitTreeState: dirty
  GoVersion: go1.17.0
`,
};
