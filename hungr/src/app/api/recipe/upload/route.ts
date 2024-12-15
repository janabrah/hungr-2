import { put } from "@vercel/blob";
import { NextResponse } from "next/server";
import { createClient } from "@supabase/supabase-js";
import { createHash } from "crypto";

const imageBypass = false;
const metadataBypass = false;

const supabase = createClient(
  process.env.NEXT_PUBLIC_SUPABASE_URL!,
  process.env.SUPABASE_SERVICE_ROLE_KEY!
);

export async function POST(request: Request): Promise<NextResponse> {
  const { searchParams } = new URL(request.url);
  const filename = searchParams.get("filename");
  console.log("filename", filename);
  const tagString = searchParams.get("tagString");
  console.log("tagString", tagString);
  if (!filename) {
    throw "filename is required";
  }
  const image = await sendImage(filename, request.body);
  const imageBlob = await image.json();
  console.log("imageblob", imageBlob, JSON.stringify(imageBlob));
  const metadata = await sendMetadata(filename, tagString, imageBlob.url);
  return NextResponse.json({ image, metadata });
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
  console.log("sending image");
  if (imageBypass) {
    return NextResponse.json(null);
  }
  // ⚠️ The below code is for App Router Route Handlers only
  if (!requestBody) {
    throw "filename and request body is required";
  }
  const blob = await put(filename, requestBody, {
    access: "public",
  });
  console.log("blob", blob);

  return NextResponse.json(blob);
}

async function sendMetadata(
  filename: string,
  tags: string | null,
  imageUrl: string
): Promise<NextResponse> {
  // ⚠️ The below code is for App Router Route Handlers only
  if (!tags || metadataBypass) {
    return NextResponse.json(null);
  }
  console.log("in sendmetadata, tags", tags);

  try {
    // Insert file into 'files' table
    console.log("setting file in files");
    const { data: file, error: fileError } = await supabase
      .from("files")
      .insert([{ filename, url: imageUrl, user_id: 1, tag_string: tags }])
      .select()
      .single();

    console.log("set file");

    if (fileError) {
      throw new Error(fileError.message);
    }

    // Insert tags into 'tags' table
    const tagInserts = tags
      .split(", ")
      .map((tag: string) => ({ id: createID(tag), name: tag }));
    console.log("setting tags in tags", tagInserts);
    const { data: insertedTags, error: tagError } = await supabase
      .from("tags")
      .upsert(tagInserts, { onConflict: "id" })
      .select();
    console.log("set tags, insertedTags", insertedTags);

    if (tagError) {
      throw new Error(tagError.message);
    }

    // Link tags to the file in 'file_tags' table
    console.log("setting fileTagLinks in file_tags");
    const fileTagLinks = insertedTags.map((tag) => ({
      file_id: file.id,
      tag_id: tag.id,
    }));
    console.log("set fileTagLinks", fileTagLinks);

    // Upload tags to 'file_tags' table
    const { error: linkError } = await supabase
      .from("file_tags")
      .insert(fileTagLinks);
    console.log("linkError", linkError);

    if (linkError) {
      throw new Error(linkError.message);
    }

    return NextResponse.json({ success: true, file, tags: insertedTags });
  } catch (error) {
    if (error instanceof Error) {
      console.error("Error uploading recipe:", error.message);
      return NextResponse.json({ error: error.message }, { status: 500 });
    } else {
      console.error("Unknown error:", error);
      return NextResponse.json(
        { error: "An unknown error occurred" },
        { status: 500 }
      );
    }
  }
}

function createID(str: string): number {
  const hash = createHash("sha256").update(str).digest("hex");
  console.log(hash);
  console.log(hash.slice(0, 8));
  return parseInt(hash.slice(0, 8), 16);
}
