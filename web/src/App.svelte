<script lang="ts">
  import Layout from "./lib/layout.svelte";
  import RecipePage from "./lib/recipe_page/recipe_page.svelte";
  
  function handleNavMsg(event: { detail: { text: any; }; }) {
    switch (event.detail.text) {
    case 'recipe':
      innerState = 'recipe'
      break;
    case 'settings':
      innerState = 'settings'
      break;
    default:
      break;
    }
  }

  let innerState = 'recipe';
</script>

<main>
<Layout on:navclick={handleNavMsg}>
  {#if innerState == 'recipe'}
    <RecipePage />
  {:else}
    <div class="text-xl text-red-700 flex justify-center m-10">
      Invalid page state reached.
      <br>
      (src/App.svelte)
    </div>
  {/if}
</Layout>
</main>
 
