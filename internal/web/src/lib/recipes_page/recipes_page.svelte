<script lang="ts">
  import RecipeItem from "./recipe_item.svelte";

  let recipes = [
    { id: 1, title: 'easy pancakes'}, 
    { id: 2, title: 'waffles'}, 
    { id: 3, title: 'french toast'}, 
    { id: 4, title: 'breakfast burrito'}, 
  ];

  function handleItmMsg(event: { detail: { tag: string, msg: string } }) {
    let id = recipes.findIndex((v=>v.title==event.detail.msg));
    recipes.splice(id, 1);
    recipes = recipes;
  }
</script>

<div class="mx-5 my-2 rounded-box flex-col">
  {#if recipes.length > 0}
  <!-- Search bar -->
  <div class="flex justify-center my-2">
    <input type="text" placeholder="Search Recipes" 
                     class="input input-bordered input-sm w-full max-w-sm"/>
  </div>

  <!-- Recipe List -->
  <div class="flex justify-center">
    <ul class="[&>*]:text-neutral-content bg-neutral rounded-box max-w-md w-full p-2">
      {#each recipes as recipe (recipe.id)}
        <li class="my-2"><RecipeItem title={recipe.title} on:msg={handleItmMsg}/></li>
      {/each}
    </ul>
  </div>
  {:else}
      <p class="m-10 flex justify-center">No recipes :(</p>
  {/if}
</div>
