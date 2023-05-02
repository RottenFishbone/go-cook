<script lang="ts">
  import {onMount} from 'svelte';
  import RecipeItem from "./recipe_item.svelte";

  async function fetchRecipeList() {
    const resp = await fetch('/api/0/recipes/allNames');
    if (resp.ok) {
      return resp.json();
    } else {
      throw new Error("erm");
    }
  }

  let recipes: [string];
  onMount(async () => {
    recipes = await fetchRecipeList();
  })

  function handleItmMsg(event: { detail: { tag: string, msg: string } }) {
    let msg = event.detail.msg;
    switch (event.detail.tag) {
      case 'delete':
        fetch("/api/0/recipes/byName?name="+msg, {
          method: 'DELETE',
        }).then(resp => {
          if (resp.ok){
            let id = recipes.findIndex((v=>v==msg));
            recipes.splice(id, 1);
            recipes = recipes;
          } else {
            // TODO convert to a toast
            console.log("Recipe delete rejected: "+resp.status + " " + resp.statusText);
          }
        }).catch((err) => {
          // TODO convert to a toast
          alert("Failed to connect to API server: " + err)
        });
        break;
      case 'edit':
        //TODO
        break;
      case 'view':
        //TODO
        break;
      default:
        console.log("Unknown message recieved from recipe_item component.")
    }
  }
</script>

<div class="mx-5 my-2 rounded-box flex-col">
  {#if recipes}
  {#if recipes.length > 0}
  <!-- Search bar -->
  <div class="flex justify-center my-2">
    <input type="text" placeholder="Search Recipes" 
                     class="input input-bordered input-sm w-full max-w-sm"/>
  </div>
  <!-- Recipe List -->
  <div class="flex justify-center">
    <ul class="[&>*]:text-neutral-content bg-neutral rounded-box max-w-md w-full p-2">
        {#each recipes as recipe (recipe)}
          <li class="my-2"><RecipeItem recipeName={recipe} on:msg={handleItmMsg}/></li>
        {/each}
    </ul>
  </div>
  {:else}
      <p class="m-10 flex justify-center">No recipes :(</p>
  {/if}
  {/if}
</div>
