export interface Article {
  title: string;
  body: string;
}

export interface PostedArticle extends Article {
  id: number;
}

export function isPostedArticle(json: any): json is PostedArticle {
  return (json as PostedArticle).id !== undefined;
}
