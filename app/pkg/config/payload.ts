import { Configuration } from '../../../../plato-client/src/lib/clients/plato/runtime';

let config: Configuration;

export function setConfigToken(token: string) {
    config = new Configuration({
        accessToken: token
    })
}

export const  getConfig = () : Configuration | undefined => {
    if (config) {
        return config
    }
    return undefined
}