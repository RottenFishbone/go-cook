<script lang="ts">
  import {createEventDispatcher} from 'svelte';

  import Icon from "./icon.svelte";

  const dispatch = createEventDispatcher();
	function clickEvent(msg: string) {
		dispatch('navclick', {
			text: msg
    });
  }

</script>

<div class="navbar">
  <!-- Home button -->
  <div class="flex-1">
    <button class="btn btn-ghost normal-case text-xl" on:click="{_=>clickEvent('recipe')}">
      <div class="text-primary">Go</div>
      Cook
    </button>
  </div>

  <!-- Nav buttons (recipes/settings..) -->
  <div class="flex-none">
    <ul class="menu menu-horizontal px-1">
      <li>
        <button class="btn btn-ghost normal-case" on:click="{_=>clickEvent('recipe')}">
          Recipes
        </button>
      </li>
      <li>
        <button class="hidden sm:flex btn btn-ghost btn-disabled normal-case">Settings</button>
        <button tabindex="0" class="btn btn-ghost btn-circle btn-disabled sm:hidden">
          <Icon name="gear" class="h-5 w-5" width="24" height="24" />
        </button>
      </li>
    </ul>
  </div>

  <div class="divider sm:divider-horizontal"/>

  <!-- Shopping Cart --->
  <div class="flex-none">
  </div>
</div>

  <slot/>

<footer class="footer footer-center p-10 text-neutral-content">
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
