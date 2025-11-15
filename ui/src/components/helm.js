import React from "react";
import Editor from "react-simple-code-editor";
import Preview from "./preview";
//import debounce from "lodash.debounce";
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
      rawRenderValues: "",
      rawChart: defaults.rawChart,
      rawRelease: defaults.rawRelease,
      rawCapabilities: defaults.rawCapabilities,
      renderedTemplate: "",
      renderedTemplateFiles: [],
      renderError: "",
    };
  }

  componentDidMount() {
    this.updateHelmRender();
  }

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
          this.setState({
              rawChart: data.chart,
              rawRelease: data.release,
              rawValues: data.values,
              rawRenderValues: data.renderValues,
              renderedTemplate: data.preview,
              renderedTemplateFiles: data.previewFiles,
              renderError: "",
          })
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
        </div>
        <div className="container">
          <div className="input">
            <div className="input__values">
              <Tabs>
                <TabList>
                  <Tab>Chart</Tab>
                  <Tab>Values</Tab>
                  <Tab>Render Values</Tab>
                  <Tab>Release</Tab>
                </TabList>
                  <TabPanel>
                      <Editor
                          value={this.state.rawChart}
                          highlight={highlighter}
                          padding={padding}
                          style={style}
                          placeholder="Insert the contents of your Chart.yaml file (e.g name)"
                          className="input__chart__editor editor"
                      />
                  </TabPanel>
                <TabPanel>
                  <Editor
                    value={this.state.rawValues}
                    highlight={highlighter}
                    padding={padding}
                    style={style}
                    placeholder="Insert the contents of your values files"
                    className="input__values__editor editor"
                  />
                </TabPanel>

                  <TabPanel>
                    <Editor
                        value={this.state.rawRenderValues}
                        highlight={highlighter}
                        padding={padding}
                        style={style}
                        placeholder="Insert the contents of your values files"
                        className="input__values__editor editor"
                    />
                </TabPanel>
                <TabPanel>
                  <Editor
                    value={this.state.rawRelease}
                    highlight={highlighter}
                    padding={padding}
                    style={style}
                    placeholder="Insert release values (e.g Name, Namespace, etc...)"
                    className="input__release__editor editor"
                  />
                </TabPanel>
              </Tabs>
            </div>
          </div>
          <div className="preview">
            <Tabs>
              <TabList>
                  { this.state.renderedTemplateFiles.map(file => <Tab key={`l-${file.filename}`}>{file.filename}</Tab>) }
              </TabList>
                { this.state.renderedTemplateFiles.map(file => <TabPanel key={`p-${file.filename}`}>
                    <Preview
                        value={file.preview}
                        highlight={highlighter}
                        padding={padding}
                        style={style}
                        className="preview__highlighted"
                    />
                </TabPanel>) }
            </Tabs>

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
  rawChart: ``,
  rawRelease: ``,
};
