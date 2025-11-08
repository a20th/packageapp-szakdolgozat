<script lang="ts">
    import {
        Button,
        Alert,
        FormGroup,
        Input,
        ButtonGroup,
    } from "@sveltestrap/sveltestrap";
    import * as m from "$lib/paraglide/messages";
    import type { PageProps } from "./$types";
    import { enhance } from "$app/forms";
    import { getLocale } from "$lib/paraglide/runtime.js";
    let { data, form } = $props();
    let isLogin = $state(true);
    let submitInProgress = $state(false);
    let errors = { email: "" };
</script>

{#if data.loggedIn}
    <main class="container">
        <p>user logged in</p>
    </main>
{:else}
    <main class="container-md login">
        {#if data.success === true}
            <Alert dismissible color="success">
                <h4 class="alert-heading text-capitalize">Success</h4>
                {data.message}
            </Alert>
        {:else if data.success === false}
        <Alert dismissible color="danger">
                <h4 class="alert-heading text-capitalize">Error</h4>
                {data.message}
            </Alert>
        {/if}

        {#if form?.error}
            <Alert dismissible color="danger">
                <h4 class="alert-heading text-capitalize">Error</h4>
                {form.error}
            </Alert>
        {:else if form?.success}
            <Alert dismissible color="success">
                <h4 class="alert-heading text-capitalize">Success</h4>
                {form.success}
            </Alert>
        {/if}

        <ButtonGroup class="w-100">
            <Button
                on:click={() => {
                    isLogin = true;
                    form = null;
                }}
                outline
                active={isLogin}>{m.app_navbar_login()}</Button
            >
            <Button
                on:click={() => {
                    isLogin = false;
                    form = null;
                }}
                outline
                active={!isLogin}>{m.form_register_submit()}</Button
            >
        </ButtonGroup>
        {#if isLogin}
            <form
                action="?/login"
                method="POST"
                use:enhance={() => {
                    submitInProgress = true;

                    return async ({ update }) => {
                        await update();
                        submitInProgress = false;
                    };
                }}
                data-sveltekit-reload
            >
                <FormGroup floating label="E-mail">
                    <Input
                        autocomplete="email"
                        disabled={submitInProgress}
                        type="email"
                        name="email"
                        feedback={errors.email}
                        invalid={errors.email != ""}
                        required
                    />
                </FormGroup>
                <FormGroup floating label="Password">
                    <Input
                        autocomplete="current-password"
                        disabled={submitInProgress}
                        type="password"
                        name="password"
                        required
                    />
                </FormGroup>
                <Button
                    block
                    disabled={submitInProgress}
                    data-sveltekit-reload
                    type="submit">{m.app_navbar_login()}</Button
                >
            </form>
        {:else}
            <form
                action="/app?/register"
                method="POST"
                use:enhance={() => {
                    submitInProgress = true;

                    return async ({ update }) => {
                        await update();
                        submitInProgress = false;
                    };
                }}
            >
                <span
                    ><span class="text-danger">*</span> - Required fields to fill</span
                >
                <FormGroup floating>
                    <Input
                        disabled={submitInProgress}
                        autocomplete="email"
                        type="email"
                        name="email"
                        required
                    />
                    <div slot="label">
                        E-mail <span class="text-danger">*</span>
                    </div>
                </FormGroup>
                <FormGroup floating>
                    <Input
                        autocomplete="name"
                        name="fullname"
                        disabled={submitInProgress}
                        required
                    />
                    <div slot="label">
                        {m.form_register_full_name()}
                        <span class="text-danger">*</span>
                    </div>
                </FormGroup>
                <FormGroup floating>
                    <Input
                        disabled={submitInProgress}
                        autocomplete="tel"
                        placeholder="Enter a value"
                        name="phone"
                        required
                    />
                    <div slot="label">
                        {m.form_register_phone_number()}
                        <span class="text-danger">*</span>
                    </div>
                </FormGroup>
                <FormGroup floating>
                    <Input
                        disabled={submitInProgress}
                        autocomplete="new-password"
                        type="password"
                        placeholder="Enter a value"
                        name="password"
                        id="password"
                        required
                    />
                    <div slot="label">
                        {m.form_login_password()}
                        <span class="text-danger">*</span>
                    </div>
                </FormGroup>
                <span class="text-wrap m-0">Jelszónak tartalmaznia kell legalább egy betűt, egy számot és 8-32 karakter között kell lennie</span>         
                <FormGroup floating>
                    <Input
                        disabled={submitInProgress}
                        autocomplete="new-password"
                        type="password"
                        name="password-confirm"
                        placeholder="Enter a value"
                        required
                    />
                    <div slot="label">
                        {m.form_register_password_again()}
                        <span class="text-danger">*</span>
                    </div>
                </FormGroup>
                <input type="hidden" name="lang" value={getLocale()} />
                <Button block disabled={submitInProgress} type="submit"
                    >{m.form_register_submit()}</Button
                >
            </form>
        {/if}
    </main>
{/if}
