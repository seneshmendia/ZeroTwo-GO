module.exports = {
  apps: [
    {
      name: "your-go-app",
      script: "./path/to/your/go/app",
      watch: true,
      ignore_watch: ["node_modules", "logs"],
    },
  ],
};
