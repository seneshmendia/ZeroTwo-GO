module.exports = {
  apps: [
    {
      name: "ZeroTwoGo",
      script: "./main.go",
      watch: true,
      ignore_watch: ["node_modules", "logs"],
    },
  ],
};
