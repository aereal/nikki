import { type CodegenConfig } from '@graphql-codegen/cli'
import { type ClientPresetConfig } from '@graphql-codegen/client-preset'

const config: CodegenConfig = {
  ignoreNoDocuments: true,
  schema: '../schema.gql',
  documents: ['./src/**/*.vue'],
  hooks: { afterAllFileWrite: ['prettier -w'] },
  config: {
    strictScalars: true,
    defaultScalarType: 'unknown',
    enumsAsTypes: true,
    useTypeImports: true,
    immutableTypes: true,
    skipTypename: true,
    scalars: {
      DateTime: {
        input: 'string',
        output: 'Date',
      },
    },
  },
  generates: {
    'src/graphql/': {
      preset: 'client',
      presetConfig: {
        fragmentMasking: { unmaskFunctionName: 'getFragmentData' },
      } satisfies ClientPresetConfig,
    },
  },
}
export default config
