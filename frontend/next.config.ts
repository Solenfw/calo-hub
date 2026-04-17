import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  allowedDevOrigins: ["http://localhost:3000", "172.27.2.166"],
  output: "standalone"
};

export default nextConfig;
