import { API_ORIGIN } from "../endpoints";
import { Article, isPostedArticle , PostedArticle } from "../models/article";

const postArticleMutation = `
mutation post($article: ArticleInputType!) {
  postArticle(article: $article) {
    id
    title
    body
  }
}
`;

const updateArticleMutation = `
mutation update($articleId: ID!, $article: ArticleUpdateInputType!) {
  updateArticle(articleId: $articleId, article: $article) {
    id
    title
    body
  }
}
`;

export const postArticle = (token: string, article: Article): Promise<PostedArticle> => {
  const variables = {
    article: {
      body: article.body,
      title: article.title,
    },
  };
  const req = window.fetch(`${API_ORIGIN}/graphql`, {
    body: JSON.stringify({
      query: postArticleMutation,
      variables,
    }),
    credentials: "same-origin",
    headers: {
      "authorization": `bearer ${token}`,
      "content-type": "application/json",
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

export const updateArticle = (token: string, article: PostedArticle): Promise<PostedArticle> => {
  const variables = {
    article: {
      body: article.body,
      title: article.title,
    },
    articleId: article.id,
  };
  const req = window.fetch(`${API_ORIGIN}/graphql`, {
    body: JSON.stringify({
      query: updateArticleMutation,
      variables,
    }),
    credentials: "same-origin",
    headers: {
      "authorization": `bearer ${token}`,
      "content-type": "application/json",
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
