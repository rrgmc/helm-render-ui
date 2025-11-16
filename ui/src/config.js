const dev = {
    api: "http://localhost:17821",
};

const prod = {
    api: "",
};

console.log("ENV: " + process.env.NODE_ENV)
const config = process.env.NODE_ENV === "development" ? dev: prod;

export default config;
