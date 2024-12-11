import { put } from "@vercel/blob";
import { NextResponse } from "next/server";

export async function POST(request: Request): Promise<NextResponse> {
  console.log(request);
  const { searchParams } = new URL(request.url);
  console.log(searchParams);
  const filename = searchParams.get("filename");
  console.log("body is", request.body);
  // ⚠️ The below code is for App Router Route Handlers only
  if (!filename || !request.body) {
    throw "filename and request body is required";
  }
  const blob = await put(filename, request.body, {
    access: "public",
  });

  // Here's the code for Pages API Routes:
  // const blob = await put(filename, request, {
  //   access: 'public',
  // });

  return NextResponse.json(blob);
}

// The next lines are required for Pages API Routes only
// export const config = {
//   api: {
//     bodyParser: false,
//   },
// };
