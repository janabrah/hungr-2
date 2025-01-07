import { put, PutBlobResult } from "@vercel/blob";
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
  const formData = await request.formData();
  console.log("formData", formData);
  const files = formData.getAll("file");
  console.log("files", files);
  const filename = searchParams.get("filename");
  console.log("filename", filename);
  const tagString = searchParams.get("tagString");
  console.log("tagString", tagString);
  if (!filename) {
    throw "filename is required";
  }
  const images = await sendImages(filename, files);
  console.log("images", images, JSON.stringify(images));
  const imagesJson = await images.json();
  console.log("imagesJson", imagesJson, JSON.stringify(imagesJson));
  const imageBlobs = imagesJson.map((image: PutBlobResult) => image.url);
  console.log("imageblob", imageBlobs, JSON.stringify(imageBlobs));
  const metadata = await sendMetadata(filename, tagString, imageBlobs);
  return NextResponse.json({ images, metadata });
}

export async function GET(request: Request): Promise<NextResponse> {
  const { searchParams } = new URL(request.url);
  if (searchParams.get("type") == "image") {
    return getImage(searchParams);
  } else {
    return getImageOptions(searchParams);
  }
}

async function getImage(searchParams: URLSearchParams): Promise<NextResponse> {
  // copilot generated, is wrong
  const url = searchParams.get("imageUrl");
  if (!url) {
    throw "url is required";
  }
  const res = await fetch(url);
  console.log(res);
  return NextResponse.json(res);
}

async function getImageOptions(
  searchParams: URLSearchParams
): Promise<NextResponse> {
  // copilot generated, is wrong
  const { data: recipeData, error: recipeError } = await supabase
    .from("recipes")
    .select("id, filename, tag_string, created_at")
    .eq("user_id", searchParams.get("user_id"))
    .order("created_at", { ascending: false })
    .range(0, 100);
  if (recipeError) {
    console.log("error was: " + recipeError.message);
    throw recipeError.message;
  }
  console.log("recipeData was: " + JSON.stringify(recipeData));
  const recipeIds = recipeData.map((recipe) => recipe.id);
  const { data: mappingData, error: mappingError } = await supabase
    .from("file_recipes_temp")
    .select("file_id, recipe_id")
    .in("recipe_id", recipeIds)
    .range(0, 10000);
  if (mappingError) {
    console.log("error was: " + mappingError.message);
    throw mappingError.message;
  }
  console.log("mappingData was: " + JSON.stringify(mappingData));
  const fileIds = mappingData.map((mapping) => mapping.file_id);
  const { data: fileData, error: fileError } = await supabase
    .from("files_temp")
    .select("id, url")
    .in("id", fileIds)
    .range(0, 10000);
  if (fileError) {
    console.log("error was: " + fileError.message);
    throw fileError.message;
  }
  console.log("fileData was: " + JSON.stringify(fileData));
  return NextResponse.json({ recipeData, fileData, mappingData });
}

async function sendImages(
  filename: string,
  files: FormDataEntryValue[]
): Promise<NextResponse> {
  console.log("sending image");
  if (imageBypass) {
    console.log("bypassing image");
    return NextResponse.json(null);
  } else {
    console.log("not bypassing image");
  }
  console.log("files is", files, JSON.stringify(files));
  const blobs = [];
  let pageNum = 0;
  for (const file of files) {
    pageNum++;
    const blob = await put(filename + pageNum, file, {
      access: "public",
    });
    console.log("blob", blob);
    blobs.push(blob);
  }

  return NextResponse.json(blobs);
}

async function sendMetadata(
  filename: string,
  tags: string | null,
  imageUrls: string[]
): Promise<NextResponse> {
  // ⚠️ The below code is for App Router Route Handlers only
  if (!tags || metadataBypass) {
    return NextResponse.json(null);
  }
  console.log("in sendmetadata, tags", tags);

  try {
    console.log("sending recipe: " + filename + " with tags: " + tags);
    const recipe = await writeToTable("recipes", {
      filename,
      user_id: 1,
      tag_string: tags,
    });
    console.log("sending urls + " + imageUrls);
    console.log("sending urls (JSON) + " + JSON.stringify(imageUrls));

    const files = [];
    for (const imageUrl of imageUrls) {
      const file = await writeToTable("files_temp", {
        url: imageUrl,
        image: true,
      });
      files.push(file);
    }

    let pageNum = 0;
    const fileIdPayload = files.map((file) => ({
      file_id: file.id,
      recipe_id: recipe.id,
      page_number: pageNum++,
    }));

    const recipeFileLinks = await writeToTable(
      "file_recipes_temp",
      fileIdPayload,
      false
    );

    console.log("recipeFileLinks is: " + JSON.stringify(recipeFileLinks));

    const tagPayload = tags
      .split(", ")
      .map((tag: string) => ({ id: createID(tag), name: tag }));

    const insertedTags = await writeToTable("tags", tagPayload, false, true);

    if (!Array.isArray(insertedTags)) {
      throw new Error("insertedTags is not an array");
    }

    const fileTagLinks = insertedTags.map((tag) => ({
      recipe_id: recipe.id,
      tag_id: tag.id,
    }));

    // Upload tags to 'recipe_tags' table
    const tagLinksResult = writeToTable("recipe_tags", fileTagLinks, false);

    console.log("tagLinksResult is: " + JSON.stringify(tagLinksResult));

    return NextResponse.json({ success: true, recipe, tags: insertedTags });
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

async function writeToTable(
  table: string,
  payload: object,
  selectSingle: boolean = true,
  upsert: boolean = false
) {
  console.log("setting: " + JSON.stringify(payload));
  let result;
  let err;
  if (upsert) {
    const { data: myData, error: myError } = await supabase
      .from(table)
      .upsert(payload, { onConflict: "id" })
      .select();

    result = myData;
    err = myError;
  } else {
    if (selectSingle) {
      const { data: myData, error: myError } = await supabase
        .from(table)
        .insert([payload])
        .select()
        .single();

      result = myData;
      err = myError;
    } else {
      const { data: myData, error: myError } = await supabase
        .from(table)
        .insert(payload)
        .select("*");

      result = myData;
      err = myError;
    }
  }

  if (err) {
    throw new Error(err.message);
  }
  console.log("set data and got: " + JSON.stringify(result));
  return result;
}
