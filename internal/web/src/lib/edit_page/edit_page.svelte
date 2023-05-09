<script lang="ts">
  import { onMount } from 'svelte';
  
  import { apiRoot } from '../common';
  import RecipePage from '../recipe_page/recipe_page.svelte';

  export let recipeName: string = '';
  
  let recipeText: string = '';
  $: textArea = recipeText

  let titleInput = recipeName;
  $: titleChanged = titleInput != recipeName;
	$: titleIllegal = titleInput == "";
	
	let titleInputClass = '';
	$: if (titleIllegal) { 
		titleInputClass = 'input-error'; 
	} else if (titleChanged && !newRecipeMode) {
		titleInputClass = 'input-info';
	} else {
		titleInputClass = '';
	}


  let failedLoad = false;     // Flag: inability to fetch recipe file
  let mounted = false;        // Flag: component has been loaded for a time
  $: newRecipeMode = false;  	// Flag: no recipe will be loaded on mount
	let recipeFetched = false; 	// Flag: recipe has been loaded (even if empty)

  // Code run once this component is mounted completely
  onMount(async ()=>{
    setTimeout(()=>{
      mounted = true;
    }, 200);

    // Attempt to fetch the relevant file based on the passed recipeName
    if (recipeName != ''){
      try {
				recipeText = await fetchRecipeFile(recipeName);
				recipeFetched = true;
      } catch (err) {
        failedLoad = true;
        throw err;
      }
    } else {
      // Otherwise we are editing an non-existant recipe 
			newRecipeMode = true;
			recipeFetched = true;
    }
  });

  // Fetches the recipe's raw file from the API server.
  async function fetchRecipeFile(name: string) {
    const resp = await fetch(`${apiRoot}/recipes/byName?name=${name}&raw=true`);
    if (resp.ok) {
      return await resp.text()
    } else {
      throw new Error(`Failed to fetch recipe file: ${resp.status} ${resp.statusText}`);
    }
  }

	// Tab handling
  function sourceTabClicked() {
    if (tabState == TabState.Source) { return; }
    tabState = TabState.Source;
    
  }
	// Tab handling
  function previewTabClicked(){
    if (tabState == TabState.Preview) { return; }
    tabState = TabState.Preview;
  }
	enum TabState {
    Source,
    Preview,
  }
	let tabState = TabState.Source;


	// Handler for saving changes made since page load
  // newFile: boolean - Determines if the changes should create a new recipe 
  //                    update an existing one
	function saveClick(newFile: boolean) {
		if (titleIllegal) {
      // Set saving to failure state
      saving = SaveState.Err;
      // Reset after a second
      setTimeout(() =>{
        saving = SaveState.Default;
			}, 2000);
			return;
		}

    saving = SaveState.Saving; // disables save button
    
    saveChanges(newFile).then(()=>{
      // Set saving to confirmation state
      saving = SaveState.Saved; 
      // Reset after a second
      setTimeout(() =>{
        saving = SaveState.Default;
      }, 1000);

      if (titleChanged) {
        // "Move" this page's state to newly created recipe
				recipeName = titleInput;
				recipeText = textArea;
				newRecipeMode = false;		// Recipe definitely exists now
      }
    }).catch(err=>{
      // Set saving to failure state
      saving = SaveState.Err;
      // Reset after a second
      setTimeout(() =>{
        saving = SaveState.Default;
      }, 2000);
      console.error(err);
    }); 
    return;
  }

  enum SaveState {
    Default,
    Saving,
    Saved,
    Err,
  }
	let saving = SaveState.Default;

	// Pushes state change to API server
  // newFile: boolean - Determines if the changes should create a new recipe 
  //                    update an existing one
	async function saveChanges(newFile: boolean) {
    let resp: Response;
    let textAreaChanged = textArea != recipeText;
		let body = textAreaChanged || newRecipeMode ? 
				(textArea == '' ? `-- im empty inside ='(` : textArea) : '';

    // Handle "Save Updates"
    if (!newFile) {
      if (!textAreaChanged && !titleChanged) { return; }
      let renameParam = titleChanged ? `&rename=${titleInput}` : '';
      let reqUrl = `${apiRoot}/recipes/?name=${recipeName}${renameParam}`;
      console.log(`req: ${reqUrl}`)
      resp = await fetch(reqUrl, {
        method: 'POST',
        body: body,
      });
    } else {
      // Handle "Save as New"
      resp = await fetch(`${apiRoot}/recipes/?name=${titleInput}`, {
        method: 'PUT',
        body: body,
      });
    }

    if (resp.ok){
      return;
    } else {
      throw new Error("Failed to save changes.");
    }
  }



  // A callback to be used alongside an svg `use:` directive
  // Resets the animations current time to 0 on mounting to DOM
  function animRestart(node: SVGSVGElement) {
    node.setCurrentTime(0);
  }

</script>
{#if recipeFetched}
  <div class="mx-auto max-w-[66rem]">
    <!-- Title -->
		<div class="form-control w-full">
			<label class={`label ${titleChanged || recipeName == '' ? '' : 'opacity-[0.01]'}`}>
				<span class="label-text-alt">
					{#if !titleIllegal && !newRecipeMode}
						recipe will be renamed on save
					{:else if titleInput == ''}
						title cannot be empty
					{:else if !newRecipeMode}
						invalid title
					{:else}
						<br/><!-- Placeholder to prevent removal from DOM -->
					{/if}
				</span>
			</label>
      <input type="text" bind:value={titleInput} class={`
						 input mb-5 w-full lower-z ${titleInputClass}`}/>
		</div>

    <!-- Tab buttons -->
    <div class="tabs tabs-boxed bg-base-100">
      <button class={`tab tab-bordered ${tabState == TabState.Source ? 'tab-active' : ''}`} 
              on:click={sourceTabClicked}>
        Source
      </button>
      <button class={`tab tab-bordered ${tabState == TabState.Preview ? 'tab-active' : ''}`} 
              on:click={previewTabClicked}>
        Preview
      </button>
    </div>
    
		<!-- Source editor (stays loaded to preserve undo history) -->
			<div class={`w-full lower-z flex flex-col ${tabState == TabState.Source ? 'h-[50vh]' : 'h-1 opacity-0'}`}>

      <textarea bind:value={textArea} class="textarea grow textarea-bordered textarea-ghost w-full"/>
      <div class="btn-group flex flex-none justify-center gap-3">
        {#if saving == SaveState.Default}
        <button class={`btn flex-auto ${titleChanged || newRecipeMode ? '' : 'hidden'}`} 
           on:click={()=>{saveClick(true);}}>
          Save As New
        </button>
        <button class={`btn btn-primary flex-auto ${!newRecipeMode ? '' : 'hidden'}`}
           on:click={()=>{saveClick(false);}}>
          Save Updates
        </button>
        {:else if saving == SaveState.Saving}
          <button class="btn btn-disabled loading flex-auto"/>
        {:else if saving == SaveState.Saved}
          <!-- Save success button -->
          <button class="btn btn-disabled flex-auto">
            <svg use:animRestart class="stroke-primary w-6 h-6" viewBox="0 0 64 64" xmlns="http://www.w3.org/2000/svg">
              <polyline class="path checkmark" stroke-width="5" fill="none" opacity="0"
                        cx="32" cy="32" points="12,32 27,45 50,21">
                <animate attributeName="opacity" from="0" to="1" begin="0.5s" dur="0.01s" fill="freeze"/>
                <animate attributeName="stroke-dasharray" from="0 300" to="300 0" begin="0.5s" dur="0.75s" fill="freeze" />
              </polyline>
              <circle id="circle" cx="32" cy="32" r="29" stroke-width="4" fill="none" >
                <animate attributeName="stroke-dasharray" from="0 1000" to="1000 0" begin="0s" dur="2s" fill="freeze" />
              </circle>
            </svg>
          </button>
        {:else if saving == SaveState.Err}
          <!-- Save failure button -->
          <!-- TODO: Toast with message -->
          <button class="btn btn-disabled bg-error flex-auto">
            <svg use:animRestart class="stroke-error-content w-6 h-6" viewBox="0 0 64 64" xmlns="http://www.w3.org/2000/svg">
              <line class="path line" stroke-width="4" fill="none"
                    x1="20" y1="20" x2="43" y2="43"/>
              <line class="path line" stroke-width="4" fill="none"
                    x1="43" y1="20" x2="20" y2="43"/>
              <circle id="circle" cx="32" cy="32" r="29" stroke-width="4" fill="none" >
                <animate attributeName="stroke-dasharray" from="0 1000" to="1000 0" begin="0s" dur="2s" fill="freeze" />
              </circle>
            </svg>
          </button>

        {/if}
      </div>
    </div>

    <!-- Preview -->
    {#if tabState == TabState.Preview }
			<RecipePage recipeName={titleInput} recipeText={textArea} previewMode={true} />
    {/if}
  </div>


<!-- Loading spinner -->
{:else if !failedLoad}
  <div class={`flex flex-col justify-center mx-auto max-w-md transition-opacity duration-750 ${mounted ? '' : 'min-h-screen opacity-0'}`}>
    <div class="text-xl mx-auto my-5">Fetching recipe source...</div>
    <div class="btn btn-circle btn-xl btn-disabled mx-auto loading btn-primary"></div>
  </div>
{:else}
  <div class="flex justify-center mx-auto max-w-md">
    <div class="my-5 bg-warning text-warning-content text-lg p-2">
      Failed to load recipe file from server.<br/>
      Try Refreshing?
    </div>
  </div>
{/if}
