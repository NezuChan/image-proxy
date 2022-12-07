import "dotenv/config";

import Crypto from "crypto";
import Fastify, { FastifyReply, FastifyRequest } from "fastify";
import { fetch, FetchResultTypes } from "@kirishima/fetch";
import Sharp, { FormatEnum } from "sharp";
import { fileTypeFromBuffer } from "file-type";

const fastifyInstance = Fastify({
    logger: {
        name: "image-proxy",
        timestamp: true,
        level: process.env.NODE_ENV === "production" ? "info" : "trace",
        formatters: {
            bindings: () => ({
                pid: "Image Proxy"
            })
        },
        transport: {
            targets: [
                { target: "pino-pretty", level: process.env.NODE_ENV === "production" ? "info" : "trace", options: { translateTime: "SYS:yyyy-mm-dd HH:MM:ss.l o" } }
            ]
        }
    },
    maxParamLength: Number.MAX_SAFE_INTEGER
});

fastifyInstance.get("/:size/:key", {
    schema: {
        params: {
            type: "object",
            properties: {
                size: { type: "string" },
                key: { type: "string" }
            }
        },
        querystring: {
            type: "object",
            properties: {
                format: { type: "string" }
            }
        }
    }
}, async (request: FastifyRequest<{ Params: { size: string; key: string }; Querystring?: { format?: keyof FormatEnum } }>, reply: FastifyReply) => {
    const { size, key } = request.params;
    const [width, height] = size.split("x").map(s => parseInt(s));
    if (width > parseInt(process.env.MAX_WITDH ?? "4096") || height > parseInt(process.env.MAX_HEIGHT ?? "4096")) throw new Error("Image too large");

    const decipher = Crypto.createDecipheriv("aes-256-cbc", process.env.KEY!, process.env.IV!);
    const decrypted = decipher.update(key, "hex");
    const decryptedString = Buffer.concat([decrypted, decipher.final()]).toString();

    const responseBuffer = await fetch(decryptedString, { redirect: "follow" }, FetchResultTypes.Buffer);

    const imageBuffer = await Sharp(responseBuffer)
        .resize(width, height)
        .toFormat(request.query?.format ?? "png", { quality: 100 })
        .toBuffer();

    const mimeType = await fileTypeFromBuffer(imageBuffer);
    return reply.header("Content-Type", mimeType?.mime ?? "image/png").send(imageBuffer);
});

try {
    await fastifyInstance.listen({ port: parseInt(process.env.PORT ?? process.env.SERVER_PORT ?? "3000"), host: "0.0.0.0" });
} catch (e) {
    fastifyInstance.log.error(e);
    process.exit(1);
}
