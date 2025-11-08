<script lang="ts">
  import {
    Navbar,
    NavbarBrand,
    NavbarToggler,
    Collapse,
    Nav,
    NavItem,
    NavLink,
    Dropdown,
    DropdownItem,
    DropdownMenu,
    DropdownToggle,
  } from "@sveltestrap/sveltestrap";
  import * as m from "$lib/paraglide/messages";
  import { setLocale, getLocale } from "$lib/paraglide/runtime.js";

  let { data, children } = $props();

  let isNavOpen = $state(true);

  function handleNavUpdate(event: CustomEvent) {
    isNavOpen = event.detail.isOpen;
  }
</script>

<Navbar color="light" light expand="md" container="md">
  <NavbarBrand href="/admin">MyPackage Admin</NavbarBrand>
  <NavbarToggler on:click={() => (isNavOpen = !isNavOpen)} />
  <Collapse isOpen={isNavOpen} navbar expand="md" on:update={handleNavUpdate}>
    <Nav class="ms-auto" navbar data-sveltekit-reload>
      <NavItem>
        <NavLink
          class="pe-1 {getLocale() == 'hu' ? 'active' : ''}"
          on:click={() => setLocale("hu")}>HU</NavLink
        >
      </NavItem>
      <NavItem>
        <NavLink
          class={getLocale() == "en" ? "active" : ""}
          on:click={() => setLocale("en")}>EN</NavLink
        >
      </NavItem>
      <NavItem>
        <NavLink href="/admin/status">{m.app_navbar_tracking()}</NavLink>
      </NavItem>
      <NavItem>
        <NavLink href="/admin/order">{m.app_navbar_order()}</NavLink>
      </NavItem>
      <NavItem>
        <NavLink href="/admin/package">{m.app_navbar_pricing()}</NavLink>
      </NavItem>
      <NavItem>
        <NavLink href="/admin/users">Users</NavLink>
      </NavItem>
      {#if data.loggedIn}
        <NavItem>
          <form method="POST" action="/admin?/logout">
            <button type="submit" class="nav-link">{m.app_navbar_logout()}</button>
          </form>
        </NavItem>
      {:else}
        <NavItem>
          <NavLink href="/admin">{m.app_navbar_login()}</NavLink>
        </NavItem>
      {/if}
    </Nav>
  </Collapse>
</Navbar>

{@render children()}
