<script lang="ts">
  import { State } from "./common";
  import {createEventDispatcher} from 'svelte';
  
  const dispatch = createEventDispatcher();
	function clickEvent(state: State) {
		dispatch('msg', {
      tag: state,
      msg: ''
    });
  }

</script>

<div class="navbar">
  <!-- Home button -->
  <div class="flex-1">
    <button class="btn btn-ghost normal-case text-xl" on:click="{_=>clickEvent(State.RecipeList)}">
      <div class="text-primary">Go</div>
      Cook
    </button>
  </div>

  <!-- Nav buttons (recipes/settings..) -->
  <div class="flex-none">
    <ul class="menu menu-horizontal px-1">
      <li>
        <button class="btn btn-ghost normal-case" on:click="{_=>clickEvent(State.RecipeList)}">
          Recipes
        </button>
      </li>
      <li>
        <button class="flex btn btn-ghost btn-disabled normal-case">Settings</button>
      </li>
    </ul>
  </div>

  <div class="divider sm:divider-horizontal"/>

  <!-- Shopping Cart --->
  <div class="flex-none">
  </div>
</div>

  <slot/>

<footer class="footer footer-center p-10">
  <div class="divider">
    <p class="text-xs">
    Copyright 2023 Jayden Dumouchel 
    <br/>
    <br/>
    <a class="link link-hover text-xs" href="https://git.sr.ht/~rottenfishbone/go-cook/blob/8fee001a786f108b5a7f786bbc5ea1b460aed6e7/LICENSE">
      MIT License
    </a>
    <br/>
    <a class="link link-hover text-xs" href="https://git.sr.ht/~rottenfishbone/go-cook">
      Source Code
    </a>
    </p>
  </div>
</footer>
