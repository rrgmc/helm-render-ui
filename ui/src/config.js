const dev = {
    api: "http://localhost:17821",
};

const prod = {
    api: "",
};

const config = process.env.NODE_ENV === "development" ? dev: prod;

export default config;
