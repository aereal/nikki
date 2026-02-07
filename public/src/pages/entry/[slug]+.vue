<script lang="ts">
import { graphql } from '@/graphql'
import { useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'

const query = graphql(`
  query GetPermalink($slug: String!) {
    article(slug: $slug) {
      slug
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
  <div v-else-if="data && data.article">
    <h1>entry: {{ data.article.slug }}</h1>
  </div>
</template>
