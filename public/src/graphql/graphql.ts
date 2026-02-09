/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core'
export type Maybe<T> = T | null
export type InputMaybe<T> = T | null | undefined
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] }
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> }
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> }
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = {
  [_ in K]?: never
}
export type Incremental<T> =
  | T
  | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never }
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string }
  String: { input: string; output: string }
  Boolean: { input: boolean; output: boolean }
  Int: { input: number; output: number }
  Float: { input: number; output: number }
  DateTime: { input: string; output: Date }
}

export type Article = {
  readonly body: Scalars['String']['output']
  readonly publishedAt: Scalars['DateTime']['output']
  readonly slug: Scalars['String']['output']
  readonly title: Scalars['String']['output']
}

export type ArticleConnection = {
  readonly nodes: ReadonlyArray<Article>
}

export type ArticleOrder = {
  readonly direction: OrderDirection
  readonly field: ArticleOrderField
}

export type ArticleOrderField = 'PUBLISHED_AT'

export type OrderDirection = 'ASC' | 'DESC'

export type Query = {
  readonly article?: Maybe<Article>
  readonly articles: ArticleConnection
}

export type QueryArticleArgs = {
  slug: Scalars['String']['input']
}

export type QueryArticlesArgs = {
  first: Scalars['Int']['input']
  order: ArticleOrder
}

export type GetPermalinkQueryVariables = Exact<{
  slug: Scalars['String']['input']
}>

export type GetPermalinkQuery = {
  readonly article?: {
    readonly slug: string
    readonly title: string
    readonly body: string
    readonly publishedAt: Date
  } | null
}

export const GetPermalinkDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'query',
      name: { kind: 'Name', value: 'GetPermalink' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'slug' } },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'String' } },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'article' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'slug' },
                value: { kind: 'Variable', name: { kind: 'Name', value: 'slug' } },
              },
            ],
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                { kind: 'Field', name: { kind: 'Name', value: 'slug' } },
                { kind: 'Field', name: { kind: 'Name', value: 'title' } },
                { kind: 'Field', name: { kind: 'Name', value: 'body' } },
                { kind: 'Field', name: { kind: 'Name', value: 'publishedAt' } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetPermalinkQuery, GetPermalinkQueryVariables>
