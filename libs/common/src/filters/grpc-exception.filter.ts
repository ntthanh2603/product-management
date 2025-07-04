// import { Catch, ArgumentsHost, HttpStatus } from "@nestjs/common";
// import { BaseExceptionFilter } from "@nestjs/core";
// import { RpcException } from "@nestjs/microservices";

// @Catch(RpcException)
// export class GrpcExceptionFilter extends BaseExceptionFilter {
//   catch(exception: RpcException, host: ArgumentsHost) {
//     const ctx = host.switchToHttp();
//     const response = ctx.getResponse();
//     const request = ctx.getRequest();

//     const error = exception.getError();

//     let status = HttpStatus.INTERNAL_SERVER_ERROR;
//     let message = "Internal server error";

//     if (typeof error === "string") {
//       message = error;
//       if (error.includes("not found")) {
//         status = HttpStatus.NOT_FOUND;
//       } else if (error.includes("validation")) {
//         status = HttpStatus.BAD_REQUEST;
//       }
//     } else if (typeof error === "object" && error !== null) {
//       status = (error as any).status || HttpStatus.INTERNAL_SERVER_ERROR;
//       message = (error as any).message || "Internal server error";
//     }

//     response.status(status).json({
//       statusCode: status,
//       message,
//       timestamp: new Date().toISOString(),
//       path: request.url,
//     });
//   }
// }

import { Catch, RpcExceptionFilter, ArgumentsHost } from "@nestjs/common";
import { Observable, throwError } from "rxjs";
import { RpcException } from "@nestjs/microservices";

@Catch(RpcException)
export class GrpcExceptionFilter implements RpcExceptionFilter<RpcException> {
  catch(exception: RpcException, host: ArgumentsHost): Observable<any> {
    const error = exception.getError();

    // Handle different types of gRPC errors
    if (typeof error === "string") {
      return throwError(
        () =>
          new RpcException({
            code: 2, // UNKNOWN
            message: error,
          })
      );
    }

    if (typeof error === "object" && error !== null) {
      return throwError(
        () =>
          new RpcException({
            code: (error as any).code || 2, // UNKNOWN
            message: (error as any).message || "Unknown error",
            details: (error as any).details || undefined,
          })
      );
    }

    return throwError(
      () =>
        new RpcException({
          code: 2, // UNKNOWN
          message: "Internal server error",
        })
    );
  }
}
