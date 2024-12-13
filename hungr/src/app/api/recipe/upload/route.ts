import { put } from "@vercel/blob";
import { NextResponse } from "next/server";
import { createClient } from "@supabase/supabase-js";

const supabase = createClient(
  process.env.NEXT_PUBLIC_SUPABASE_URL!,
  process.env.SUPABASE_SERVICE_ROLE_KEY!
);

export async function POST(request: Request): Promise<NextResponse> {
  const { searchParams } = new URL(request.url);
  const filename = searchParams.get("filename");
  const tagString = searchParams.get("tagString");
  if (!filename) {
    throw "filename is required";
  }
  const imagePromise = sendImage(filename, request.body);
  const metadataPromise = sendMetadata(filename, tagString);
  return Promise.all([imagePromise, metadataPromise]);
}

// The next lines are required for Pages API Routes only
// export const config = {
//   api: {
//     bodyParser: false,
//   },
// };

async function sendImage(
  filename: string,
  requestBody: ReadableStream<Uint8Array> | null
): Promise<NextResponse> {
  // ⚠️ The below code is for App Router Route Handlers only
  if (!requestBody) {
    throw "filename and request body is required";
  }
  const blob = await put(filename, requestBody, {
    access: "public",
  });

  // Here's the code for Pages API Routes:
  // const blob = await put(filename, request, {
  //   access: 'public',
  // });

  return NextResponse.json(blob);
}

async function sendMetadata(
  filename: string,
  tags: string | null
): Promise<NextResponse> {
  // ⚠️ The below code is for App Router Route Handlers only
  if (!tags) {
    return NextResponse.json(null);
  }

  try {
    // Step 1: Insert file into 'files' table
    const { data: file, error: fileError } = await supabase
      .from("files")
      .insert([{ filename }])
      .select()
      .single();

    if (fileError) {
      throw new Error(fileError.message);
    }

    // Step 2: Insert tags into 'tags' table (using upsert to avoid duplicates)
    const tagInserts = tags.map((tag: string) => ({ name: tag }));
    const { data: insertedTags, error: tagError } = await supabase
      .from("tags")
      .upsert(tagInserts, { onConflict: "name" })
      .select();

    if (tagError) {
      throw new Error(tagError.message);
    }

    // Step 3: Link tags to the file in 'file_tags' table
    const fileTagLinks = insertedTags.map((tag: any) => ({
      file_id: file.id,
      tag_id: tag.id,
    }));

    const { error: linkError } = await supabase
      .from("file_tags")
      .insert(fileTagLinks);

    if (linkError) {
      throw new Error(linkError.message);
    }

    return NextResponse.json({ success: true, file, tags: insertedTags });
  } catch (error: any) {
    console.error("Error uploading recipe:", error.message);
    return NextResponse.json({ error: error.message }, { status: 500 });
  }

  return NextResponse.json(blob);
}
