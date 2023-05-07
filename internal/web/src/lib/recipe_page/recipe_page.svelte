<script lang='ts'>
  import { onMount } from 'svelte';
  
  import type { Recipe, Chunk, Component } from '../common'
  import { apiRoot, noQtyName, stripRecipeName } from '../common'

  export let recipeName: string;
  export let recipeText = '';

  let recipe: Recipe | null = null;


  // Recipe components
  let ingredients: [Component];
  let cookware: [Component];
  let timers: [Component];
  let steps: [[Chunk]];

  // Hook reactivity to components
  $: ingredients = recipe ? recipe.ingredients : null;
  $: cookware = recipe ? recipe.cookware : null;
  $: timers = recipe ? recipe.timers : null;
  $: steps = recipe ? recipe.steps : null;

  async function fetchRecipeByName(name: string) {
    const resp = await fetch(`${apiRoot}/recipes/byName?name=${name}`);
    if (resp.ok){
      return resp.json()
    } else {
      throw new Error(`Failed to fetch recipe '${name}': ${resp.status} ${resp.statusText}`);
    }
  }

  async function fetchParsedRecipe(recipeText: string) {
    const resp = await fetch(`${apiRoot}/recipes/parse`, {
      method: "POST",
      body: recipeText,
    });
    if (resp.ok){
      return resp.json()
    } else {
      throw new Error(`Failed to parse recipe: ${resp.status} ${resp.statusText}`);
    }
  }

  let loadFailed = false;
  let mounted = false;

  onMount(async () => {
    // This timer will hide content briefly to avoid flashing the user with text
    setTimeout(()=>{
      mounted = true;
    }, 200);
  
    try {
      if (recipeText == '') { 
        recipe = await fetchRecipeByName(recipeName);
      } else {
        recipe = await fetchParsedRecipe(recipeText);
      }
    } catch (err) {
      loadFailed = true;
      throw err;
    }
  });

</script>

{#if recipe}
  <!-- Main Display Container -->
  <div class="max-w-md mx-auto md:max-w-2xl text-md">
    <div class="divider flex justify-center text-2xl">{stripRecipeName(recipeName)}</div>
    <!-- arrange cards horizontally on larger screens -->
    <div class="md:flex gap-10"> 
      <!-- Ingredients Card -->
      <div class="card lower-z w-full h-min mx-auto my-4">
        <div class="card-body">
          <!-- Title -->
          <div class="card-title text-lg">Ingredients</div>
          <!-- Contents -->
          <table class="table table-compact w-full">
            {#each ingredients as ingr}
              <tr>
                {#if ingr.qty !== ""}
                  <td class="whitespace-normal break-words min-w-0 text-right">
                    {ingr.qty} {ingr.unit}
                  </td>
                {:else}
                  <td class="whitespace-normal break-words min-w-0 text-right">
                    {noQtyName}
                  </td>
                {/if}
                <td class="whitespace-normal break-words min-w-0 text-left">
                  {ingr.name}
                </td>
              </tr>
            {/each}
          </table>
        </div>
      </div>
      {#if cookware.length > 0}
      <!-- Cookware Card -->
      <div class="card lower-z rounded-box w-full h-min mx-auto my-4">
        <div class="card-body">
          <!-- Title -->
          <div class="card-title text-lg">Cookware</div>
          <!-- Contents -->
          <table class="table table-compact w-full">
            {#each cookware as cw}
              <tr>
                {#if cw.qty !== ""}
                  <td class="whitespace-normal break-words min-w-0 text-right">
                    {cw.qty} {cw.unit}
                  </td>
                {:else}
                  <!-- Place holder column so the name still goes in the 2nd col -->
                  <td class="min-w-0"/>
                {/if}
                <td class="whitespace-normal break-words min-w-0 text-left">
                  {cw.name}
                </td>
              </tr>
            {/each}
          </table>
        </div>
      </div>
      {/if}
    </div>

    <!--- Steps --->
    <div class="card lower-z rounded-box w-full h-min mx-auto my-4">
      <div class="card-body">
        <!-- Title -->
        <div class="card-title">Steps</div>

        <!-- Steps List -->
        <ol class="list-decimal list-outside md:mx-5">
          {#each steps as step}
            <li class="my-5 rounded-box list-item upper-z p-5">
              {#each step as chunk}
                {#if chunk.tag === 'text'}
                  {''+chunk.data}
                {:else if chunk.tag === 'ingredient'}
                  <span class="text-primary">{chunk.data.name}</span>
                {:else if chunk.tag === 'cookware'}
                  <span class="text-accent">{chunk.data.name}</span>
                {:else if chunk.tag === 'timer'}
                  <span class="text-info">{chunk.data.qty} {chunk.data.unit}</span>
                {/if}
              {/each}
            </li>
          {/each}
        </ol>
      </div>
    </div>
  </div>


{:else if !loadFailed}
  <div class={`flex flex-col justify-center mx-auto max-w-md transition-opacity duration-750 ${mounted ? '' : 'min-h-screen opacity-0'}`}>
    <div class="text-xl mx-auto my-5">Fetching recipe...</div>
    <div class="btn btn-circle btn-xl btn-disabled mx-auto loading btn-primary"/>
    </div>
{:else}
  <div class="flex justify-center mx-auto max-w-md">
    <div class="my-5 bg-warning text-warning-content text-lg p-2">
      Failed to load Recipe from server.<br/>
      Try Refreshing?
    </div>
  </div>
{/if}