export type Recipe = {
  Title: string;
  Description: string;
  Tags: string[];
  Filename: string;
};

export type Metadata = {
  names: string[];
  details: {
    [key: string]: Recipe;
  };
};
