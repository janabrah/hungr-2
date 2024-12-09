"use client";
import { useState } from "react";
import Image from "next/image";
import metadata from "../../../public/metadataDB.json";

type Recipe = {
  Title: string;
  Description: string;
  Tags: string[];
  Filename: string;
};

type Metadata = {
  details: {
    [key: string]: Recipe;
  };
};

const typedMetadata: Metadata = metadata;

export default function ShowRecipe() {
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
  );
}
