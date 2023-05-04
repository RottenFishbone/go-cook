export const noQtyName = 'some';

// The base adress of the API server (without the final forward slash)
export let apiRoot: string;

// Possible Navigation States
export enum State {
    RecipeList,
    RecipeView,
    Settings,
}

export interface Component {
    name:   string;
    qty:    string;
    qtyVal: number;
    unit:   string;
}

export interface Chunk {
    tag: string;
    data: Component | string;
}

export interface Recipe {
    name:           string;
    metadata:       { tag: string, body: string };
    ingredients:    [Component];
    cookware:       [Component];
    timers:         [Component];
    steps:          [[Chunk]];
}

export function stripRecipeName(name: string): string {
    return name.split('/').pop().replace('_', ' ');
}

// Swap the API server address to the Go server's local address when in dev mode
if (import.meta.env.DEV !== true){
	apiRoot = '/api/0'
} else {
	apiRoot = 'http://localhost:6969/api/0'
    console.debug(`Running web app in dev, sending API requests to ${apiRoot}`);
}



