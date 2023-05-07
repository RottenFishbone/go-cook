<script lang='ts'>
  import { onMount } from 'svelte';
  
  import { State } from './lib/common'

  import Layout from './lib/layout.svelte';
  import RecipesPage from './lib/recipes_page/recipes_page.svelte';
  import RecipePage from './lib/recipe_page/recipe_page.svelte';
    import EditPage from './lib/edit_page/edit_page.svelte';
  

  let state = State.RecipeList;
  
  let currentRecipeName: string;

  function handleNavMsg(event: { detail: { tag: any, msg: any; }; }) {
    switch (event.detail.tag) {
      case State.RecipeList:
        state = State.RecipeList;
				document.body.scrollIntoView();
				break;
      case State.RecipeView:
        state = State.RecipeView;
				currentRecipeName = event.detail.msg;
				document.body.scrollIntoView();
        break;
      case State.Settings:
        // TODO implement
        console.debug('Not implemented.');
				break;
			case State.RecipeEdit:
				state = State.RecipeEdit;
				currentRecipeName = event.detail.msg;
				document.body.scrollIntoView();
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
		<RecipePage recipeName={currentRecipeName} on:msg={handleNavMsg} />
	{:else if state == State.RecipeEdit}
		<EditPage recipeName={currentRecipeName} />
  {:else}
    <div class="text-xl text-error flex justify-center m-10">
      Invalid page state reached.
    </div>
  {/if}
</Layout>
</main>
 
