<script lang="ts">
import * as m from "$lib/paraglide/messages";
import {Accordion, AccordionItem} from "@sveltestrap/sveltestrap";

</script>

<main class="container content p-3">
    <h2 class="text-center mb-3">{m.app_navbar_faq()} - Gyakran Ismételt Kérdések</h2>
    <Accordion stayOpen>
        <AccordionItem header="Kérdés #1">
            <p>Válasz #1</p>
        </AccordionItem>
        <AccordionItem header="Kérdés #2">
            <p>Válasz #2</p>
        </AccordionItem>
        <AccordionItem header="Kérdés #3">
            <p>Válasz #3</p>
        </AccordionItem>
        <AccordionItem header="Kérdés #4">
            <p>Válasz #4</p>
        </AccordionItem>
    </Accordion>
</main>
