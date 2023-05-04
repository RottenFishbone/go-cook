<script lang='ts'>
  import { State } from './lib/common'

  import Layout from './lib/layout.svelte';
  import RecipesPage from './lib/recipes_page/recipes_page.svelte';
  import RecipePage from './lib/recipe_page/recipe_page.svelte';
  

  let state = State.RecipeList;
  
  let currentRecipeName: string;

  function handleNavMsg(event: { detail: { tag: any, msg: any; }; }) {
    switch (event.detail.tag) {
      case State.RecipeList:
        state = State.RecipeList;
        break;
      case State.RecipeView:
        state = State.RecipeView;
        currentRecipeName = event.detail.msg;
        break;
      case State.Settings:
        // TODO implement
        console.debug('Not implemented.');
        break;
      default:
        
        console.error(`handleNavMessage recieved unhandled event: ${event}`);
        break;
    }
  }

</script>

<main>
<Layout on:msg={handleNavMsg}>
  {#if state == State.RecipeList}
    <RecipesPage on:msg={handleNavMsg}/>
  {:else if state == State.RecipeView}
    <RecipePage recipeName={currentRecipeName} on:msg={handleNavMsg}/>
  {:else}
    <div class="text-xl text-red-700 flex justify-center m-10">
      Invalid page state reached.
      <br>
      (src/App.svelte)
    </div>
  {/if}
</Layout>
</main>
 
