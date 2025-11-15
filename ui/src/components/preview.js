import * as React from "react";

type Props = React.ElementConfig<"div"> & {
  // Props for the component
  value: string,
  highlight: (value: string) => string | React.Node,
  padding: number | string,
  style?: {},
};

export default class Preview extends React.Component<Props> {
  render() {
    const { value, style, padding, highlight, ...rest } = this.props;
    const contentStyle = {
      paddingTop: padding,
      paddingRight: padding,
      paddingBottom: padding,
      paddingLeft: padding,
    };
    const highlighted = highlight(value);
    return (
      <div
        style={{ ...styles.container, ...style }}
        className="preview__div"
        {...rest}
      >
        <pre
          dangerouslySetInnerHTML={{ __html: highlighted }}
          style={{ ...styles.preview, ...styles.highlight, ...contentStyle }}
          className="preview__pre"
        />
      </div>
    );
  }
}

const styles = {
  container: {
    position: "relative",
    textAlign: "left",
    padding: 0,
    boxSizing: "border-box",
  },
  highlight: {
    position: "relative",
    pointerEvents: "none",
  },
  preview: {
    margin: 0,
    border: 0,
    background: "none",
    boxSizing: "inherit",
    display: "inherit",
    fontFamily: "inherit",
    fontSize: "inherit",
    fontStyle: "inherit",
    fontVariantLigatures: "inherit",
    fontWeight: "inherit",
    letterSpacing: "inherit",
    lineHeight: "inherit",
    tabSize: "inherit",
    textIndent: "inherit",
    textRendering: "inherit",
    textTransform: "inherit",
    whiteSpace: "pre-wrap",
    wordBreak: "keep-all",
    overflowWrap: "break-word",
  },
};
