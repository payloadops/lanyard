export function disambiguator(name: string, stage: string, region: string) : string {
    return `${name}-${stage}-${region}`
}