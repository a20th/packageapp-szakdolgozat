<script lang="ts">
    import {Navbar, NavbarBrand, NavbarToggler, Collapse, Nav, NavItem, NavLink, Dropdown, DropdownItem, DropdownMenu, DropdownToggle} from '@sveltestrap/sveltestrap';
    import * as m from "$lib/paraglide/messages"
    import { setLocale, getLocale } from "$lib/paraglide/runtime.js";

	let {data, children } = $props();

    let isNavOpen = $state(true); 

    function handleNavUpdate(event : CustomEvent) {
        isNavOpen = event.detail.isOpen;
    }
</script>


<Navbar color="light" light expand="md" container="md">
  <NavbarBrand href="/app">MyPackage</NavbarBrand>
  <NavbarToggler on:click={() => (isNavOpen = !isNavOpen)} />
  <Collapse isOpen={isNavOpen} navbar expand="md" on:update={handleNavUpdate}>
    <Nav class="ms-auto" navbar data-sveltekit-reload>
        <NavItem>
            <NavLink class="pe-1 {getLocale() == "hu" ? "active" : ""}" on:click={() => setLocale("hu")}>HU</NavLink>
        </NavItem>
        <NavItem>
            <NavLink class = {getLocale() == "en" ? "active" : ""} on:click={() => setLocale("en")}>EN</NavLink>
        </NavItem>
        <NavItem>
            <NavLink href="/app/order">{m.app_navbar_order()}</NavLink>
        </NavItem>
        <NavItem>
            <NavLink href="/app/pricing">{m.app_navbar_pricing()}</NavLink>
        </NavItem>
        <NavItem>
            <NavLink href="/app/tracking">{m.app_navbar_tracking()}</NavLink>
        </NavItem>
        <NavItem>
            <NavLink href="/app/faq">{m.app_navbar_faq()}</NavLink>
        </NavItem>

      {#if data.loggedIn}
      <Dropdown nav inNavbar>
        <DropdownToggle nav caret>{data.profile?.name}</DropdownToggle>
        <DropdownMenu end>
          <DropdownItem>{m.app_navbar_profile()}</DropdownItem>
          <DropdownItem href="/app/orders">Megrendelések</DropdownItem>
          <DropdownItem divider />
          <form method="POST" action="/app?/logout">
          <DropdownItem type="submit">{m.app_navbar_logout()}</DropdownItem>
        </form>
        </DropdownMenu>
      </Dropdown>
      {:else}
        <NavItem>
        <NavLink href="/app">{m.app_navbar_login()}</NavLink>
      </NavItem>
      {/if}
      
      
    </Nav>
  </Collapse>
</Navbar>

{@render children()}



