"use client";

import type { PutBlobResult } from "@vercel/blob";
import { useState, useRef } from "react";

export default function AvatarUploadPage() {
  const [fileInputs, setFileInputs] = useState([
    { id: 1, ref: useRef<HTMLInputElement>(null) },
  ]);
  const metadataRef = useRef<HTMLInputElement>(null);
  const filenameRef = useRef<HTMLInputElement>(null);
  const [imageBlobs, setImageBlobs] = useState<PutBlobResult[]>([]);
  const [metadataBlob, setMetadataBlob] = useState<PutBlobResult | null>(null);
  const [isSubmitted, setIsSubmitted] = useState<boolean>(false);
  console.log(metadataBlob, setMetadataBlob);

  const handleFileChange = (index: number) => {
    if (index === fileInputs.length - 1) {
      setFileInputs([
        ...fileInputs,
        { id: fileInputs.length + 1, ref: useRef<HTMLInputElement>(null) },
      ]);
    }
  };

  const handleSubmit = async (event: React.FormEvent) => {
    if (isSubmitted) {
      return;
    }
    event.preventDefault();
    setIsSubmitted(true);
    try {
      await sendUpload(
        fileInputs.map((input) => input.ref),
        setImageBlobs,
        metadataRef,
        setMetadataBlob,
        filenameRef
      );
    } catch (error) {
      console.error("Error during upload:", error);
    }
  };

  return (
    <>
      <h1>Upload an image of your recipe.</h1>

      <form onSubmit={handleSubmit}>
        {fileInputs.map((input, index) => (
          <input
            key={input.id}
            name={`file-${input.id}`}
            ref={input.ref}
            type="file"
            required
            onChange={() => handleFileChange(index)}
            className="block w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        ))}
        <input
          name="metadata"
          ref={metadataRef}
          type="text"
          placeholder="Enter a list of tags, separated by commas"
          className="block w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
        />
        <input
          name="filename"
          ref={filenameRef}
          type="text"
          placeholder="Enter your desired file name"
          className="block w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
        />
        <button
          type="submit"
          className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          Upload
        </button>
      </form>
      {imageBlobs.length > 0 && (
        <div>
          {imageBlobs.map((blob, index) => (
            <div key={index}>
              Blob: <a href={blob.url}>{blob.url}</a>
            </div>
          ))}
        </div>
      )}
    </>
  );
}

async function sendUpload(
  fileRefs: React.RefObject<HTMLInputElement>[],
  setImageBlobs: React.Dispatch<React.SetStateAction<PutBlobResult[]>>,
  metadataRef: React.RefObject<HTMLInputElement | null>,
  setMetedataBlob: React.Dispatch<React.SetStateAction<PutBlobResult | null>>,
  filenameRef: React.RefObject<HTMLInputElement | null>
) {
  const files = fileRefs.flatMap((ref) =>
    ref.current?.files ? Array.from(ref.current.files) : []
  );
  if (files.length === 0) {
    throw new Error("No files selected");
  }

  const tagString = metadataRef.current?.value;
  const uploadedBlobs: PutBlobResult[] = [];

  const recipe = await createRecipe(tagString, files, filenameRef);

  let pageNum = 0;
  for (const file of files) {
    const response = await uploadFile(file, tagString, filenameRef, pageNum++);

    const newBlob = (await response.json()) as PutBlobResult;
    console.log(JSON.stringify(newBlob));
    uploadedBlobs.push(newBlob);
  }
  setImageBlobs(uploadedBlobs);
}

async function uploadFile(
  file: File,
  filenameRef: React.RefObject<HTMLInputElement | null>,
  pageNum: number
) {
  let filename = file.name;
  if (filenameRef.current?.value) {
    filename = filenameRef.current.value;
  }

  const url = `/api/recipe/upload?filename=${filename}&pageNum=${pageNum}&type=file`;

  console.log("inpage, sending url: ", url);
  const response = await fetch(url, {
    method: "POST",
    body: file,
  });

  if (!response.ok) {
    throw new Error("Failed to upload file");
  }

  return response;
}

async function createRecipe(
  tagString: string | undefined,
  files: File[],
  filenameRef: React.RefObject<HTMLInputElement | null>
) {
  let filename = files[0].name;
  if (filenameRef.current?.value) {
    filename = filenameRef.current.value;
  }

  const url = `/api/recipe/upload?filename=${filename}&tagString=${tagString}&type=recipe`;

  console.log("inpage, sending url: ", url);
  const response = await fetch(url, {
    method: "POST",
  });

  if (!response.ok) {
    throw new Error("Failed to create recipe");
  }

  return response;
}
