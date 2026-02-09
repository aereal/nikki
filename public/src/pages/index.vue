<script lang="ts">
import { graphql } from '@/graphql'
import { useQuery } from '@urql/vue'

const query = graphql(`
  query ListArticles($first: Int!) {
    articles(first: $first, order: { direction: DESC, field: PUBLISHED_AT }) {
      nodes {
        title
        slug
      }
    }
  }
`)

export default {
  setup() {
    const { fetching, data, error } = useQuery({
      query,
      variables: { first: 100 },
    })
    return {
      fetching,
      error,
      data,
    }
  },
}
</script>

<template>
  <div v-if="fetching">Loading...</div>
  <div v-else-if="error">! error: {{ error.message }}</div>
  <div v-else-if="data && data.articles">
    <section :key="article.slug" v-for="article in data.articles.nodes">
      <h1>
        <RouterLink :to="`/entry/${article.slug}`">{{ article.title }}</RouterLink>
      </h1>
    </section>
  </div>
</template>
