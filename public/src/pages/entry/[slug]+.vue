<script lang="ts">
import { graphql } from '@/graphql'
import { useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'

const query = graphql(`
  query GetPermalink($slug: String!) {
    article(slug: $slug) {
      slug
      title
      body
      publishedAt
    }
  }
`)

export default {
  setup() {
    const { params } = useRoute('/entry/[slug]+')

    const { fetching, data, error } = useQuery({
      query,
      variables: { slug: params.slug.join('/') },
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
  <article v-else-if="data && data.article">
    <header>
      <h1>{{ data.article.title }}</h1>
      <p>
        <time>{{ data.article.publishedAt }}</time>
      </p>
    </header>
    <main v-html="data.article.body" />
  </article>
</template>
