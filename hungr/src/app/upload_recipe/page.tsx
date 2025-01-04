"use client";

import type { PutBlobResult } from "@vercel/blob";
import { useState, useRef } from "react";

export default function AvatarUploadPage() {
  const inputFileRef = useRef<HTMLInputElement>(null);
  const metadataRef = useRef<HTMLInputElement>(null);
  const filenameRef = useRef<HTMLInputElement>(null);
  const [imageBlob, setImageBlob] = useState<PutBlobResult | null>(null);
  const [metadataBlob, setMetadataBlob] = useState<PutBlobResult | null>(null);
  const [isSubmitted, setIsSubmitted] = useState<boolean>(false);
  console.log(metadataBlob, setMetadataBlob);

  const handleSubmit = async (event: React.FormEvent) => {
    if (isSubmitted) {
      return;
    }
    event.preventDefault();
    setIsSubmitted(true);
    try {
      await sendUpload(
        inputFileRef,
        setImageBlob,
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
        <input
          name="file"
          ref={inputFileRef}
          type="file"
          required
          className="block w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
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
      {imageBlob && (
        <div>
          Blob url: <a href={imageBlob.url}>{imageBlob.url}</a>
        </div>
      )}
    </>
  );
}

async function sendUpload(
  inputFileRef: React.RefObject<HTMLInputElement | null>,
  setImageBlob: React.Dispatch<React.SetStateAction<PutBlobResult | null>>,
  metadataRef: React.RefObject<HTMLInputElement | null>,
  setMetedataBlob: React.Dispatch<React.SetStateAction<PutBlobResult | null>>,
  filenameRef: React.RefObject<HTMLInputElement | null>
) {
  if (!inputFileRef.current?.files) {
    throw new Error("No file selected");
  }
  const file = inputFileRef.current.files[0];
  if (!metadataRef.current?.value) {
    metadataRef.current = null;
  }
  const tagString = metadataRef.current?.value;
  let filename = file.name;
  if (filenameRef.current?.value) {
    filename = filenameRef.current.value;
  }

  const url = `/api/recipe/upload?filename=${filename}&tagString=${tagString}`;

  console.log("inpage, sending url: ", url);
  const response = await fetch(url, {
    method: "POST",
    body: file,
  });
  console.log("got response, ", response);

  const newBlob = (await response.json()) as PutBlobResult;
  console.log("newBlob", JSON.stringify(newBlob));
  console.log(setMetedataBlob);
  setImageBlob(newBlob);
}
