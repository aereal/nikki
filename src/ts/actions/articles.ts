import { API_ORIGIN } from "../endpoints";
import { Article, isPostedArticle , PostedArticle } from "../models/article";
import { AuthedUser } from "../models/user";

export const postArticle = (author: AuthedUser, article: Article): Promise<PostedArticle> => {
  const req = window.fetch(`${API_ORIGIN}/articles`, {
    body: JSON.stringify({
      body: article.body,
      title: article.title,
    }),
    credentials: "same-origin",
    headers: {
      "visitor-key": author.authKey,
    },
    method: "POST",
  });
  return req
    .then((res) => res.json())
    .then((json) => {
      if (isPostedArticle(json)) {
        return json;
      } else {
        throw new Error("Invalid response");
      }
    });
};

export const updateArticle = (author: AuthedUser, article: PostedArticle): Promise<PostedArticle> => {
  const req = window.fetch(`https://api.nikki.dev/articles/${article.id}`, {
    body: JSON.stringify({
      body: article.body,
      title: article.title,
    }),
    credentials: "same-origin",
    headers: {
      "visitor-key": author.authKey,
    },
    method: "PUT",
  });
  return req
    .then((res) => res.json())
    .then((json) => {
      if (isPostedArticle(json)) {
        return json;
      } else {
        throw new Error("Invalid response");
      }
    });
};
