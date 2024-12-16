"use client";
import { useState, useEffect } from "react";
//import { NextResponse } from "next/server";
//import Image from "next/image";
//import metadata from "../../../public/images/metadataDB.json";
const USERID = 1; // Temporarily keeping this simple

/*type Recipe = {
  Title: string;
  Description: string;
  Tags: string[];
  Filename: string;
};*/

type Metadata = {
  filename: string;
  tagString: string;
  createdAt: string;
  imageUrl: string;
};

//const typedMetadata: Metadata = metadata;

export default function ShowRecipe() {
  const [data, setData] = useState<Metadata[] | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [selectedOption, setSelectedOption] = useState<string>("");
  const [image, setImage] = useState<Blob | null>(null);

  useEffect(() => {
    fetchData();
  }, []);

  const userId = USERID;
  console.log("Entering showrecipe");
  const fetchData = async () => {
    try {
      const url = `/api/recipe/upload?type=metadata&user_id=${userId}`;
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error("Failed to fetch data");
      }
      const result = await response.json();
      setData(packageData(result));
    } catch (error) {
      if (error instanceof Error) {
        setError(error.message);
      } else {
        console.error("Unknown error:", error);
        setError("An unknown error occurred and error was the wrong type");
      }
    }
  };

  const fetchImage = async (imageUrl: string) => {
    console.log("imageUrl:", imageUrl);
    try {
      const response = await fetch(imageUrl);
      console.log("Response:", response);
      if (!response.ok) {
        throw new Error("Failed to fetch image");
      }
      console.log("Response is ok:", response);
      const result = await response.blob();
      console.log("Image result:", result);
      setImage(result);
    } catch (error) {
      if (error instanceof Error) {
        setError(error.message);
      } else {
        console.error("Unknown error:", error);
        setError("An unknown error occurred and error was the wrong type");
      }
    }
  };

  const handleSubmit = () => {
    console.log("Selected option:", selectedOption);
    const image = fetchImage(selectedOption);
    console.log(image);
  };
  console.log(image);
  return (
    <div>
      <h1>Please choose the recipe you want to load.</h1>
      {error && <div>Error: {error}</div>}
      {data != null && data.length > 0 && (
        <div>
          <h2>Data:</h2>
          <select
            value={selectedOption}
            onChange={(event) => setSelectedOption(event.target.value)}
            className="block w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="" disabled>
              Select an option
            </option>
            {data.map((item, index) => (
              <option key={index} value={item.imageUrl}>
                {item.filename} - {item.tagString}
              </option>
            ))}
          </select>
          <button
            onClick={handleSubmit}
            className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 mt-4"
          >
            Submit
          </button>
        </div>
      )}
      {image && (
        <div>
          <h2>Image:</h2>
          <img src={URL.createObjectURL(image)} alt="Fetched from database" />
        </div>
      )}
    </div>
  );

  /*
  const [selectedTitle, setSelectedTitle] = useState(
    Object.keys(typedMetadata.details)[0]
  );

  const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setSelectedTitle(event.target.value);
  };

  const selectedRecipe = typedMetadata.details[selectedTitle];

  return (
    <div>
      <main>
        <div>
          <div>Select a recipe:</div>
          <select onChange={handleChange} value={selectedTitle}>
            {Object.keys(typedMetadata.details).map((title) => (
              <option key={title} value={title}>
                {title}
              </option>
            ))}
          </select>
        </div>
        {selectedRecipe && (
          <div>
            <h2>{selectedRecipe.Title}</h2>
            <p>{selectedRecipe.Description}</p>
            <Image
              src={`/${selectedRecipe.Filename}`}
              alt={selectedRecipe.Title}
              width={500}
              height={500}
              style={{ objectFit: "contain" }}
            />
          </div>
        )}
      </main>
    </div>
  );*/
}

function packageData(data: object[]): Metadata[] {
  return data.map((item) => {
    const validItem = item as {
      filename: string;
      tag_string: string;
      created_at: string;
      url: string;
    };

    if (
      typeof validItem !== "object" ||
      !validItem ||
      typeof validItem.filename !== "string" ||
      typeof validItem.tag_string !== "string" ||
      typeof validItem.created_at !== "string" ||
      typeof validItem.url !== "string"
    ) {
      throw new Error("Data is missing required fields");
    }

    return {
      filename: validItem.filename,
      tagString: validItem.tag_string,
      createdAt: validItem.created_at,
      imageUrl: validItem.url,
    };
  });
}
