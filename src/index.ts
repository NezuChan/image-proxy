import "dotenv/config";

import Fastify from "fastify";
import Crypto from "crypto";
import Jimp from "jimp";
import { FastifyReply } from "fastify/types/reply";
import { FastifyRequest } from "fastify/types/request";

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
                key: { type: "string" },
            }
        },
        querystring: {
            type: "object",
            properties: {
                format: { type: "string", enum: ["png", "jpg", "jpeg", "webp"] },
            }
        }
    }
        }, async (request: FastifyRequest<{ Params: { size: string; key: string; }, Querystring: { format: "png" | "jpeg" }}>, reply: FastifyReply) => {
            const { size, key } = request.params;
            const [width, height] = size.split("x").map((s) => parseInt(s, 10));
            if (width > (parseInt(process.env.MAX_WITDH ?? "4096")) || height > (parseInt(process.env.MAX_HEIGHT ?? "4096"))) throw new Error("Image too large");

            const decipher = Crypto.createDecipheriv("aes-256-cbc", process.env.KEY!, process.env.IV!);
            const decrypted = decipher.update(key, "hex");
            const decryptedString = Buffer.concat([decrypted, decipher.final()]).toString();

            const image = await Jimp.read(decryptedString);
            image.resize(width, height, Jimp.RESIZE_HERMITE);
            image.quality(100);

            const buffer = await image.getBufferAsync(request.query.format === "jpeg" ? Jimp.MIME_JPEG : Jimp.MIME_PNG);

            reply.header("Content-Type", request.query.format === "jpeg" ? Jimp.MIME_JPEG : Jimp.MIME_PNG);
            return buffer;
        }
);

try {
    await fastifyInstance.listen({ port: parseInt(process.env.PORT ?? process.env.SERVER_PORT ?? "3000"), host: "0.0.0.0" });
} catch (e) {
    fastifyInstance.log.error(e);
    process.exit(1);
}
