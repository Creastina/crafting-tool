use crafting::entities::Todo;
use leptos::prelude::*;
use leptos_meta::*;
use leptos_router::components::*;
use leptos_router::path;
use web_sys::HtmlDialogElement;

#[cfg(feature = "ssr")]
pub type DbConnection = actix_web::web::Data<sea_orm::DatabaseConnection>;

#[server]
pub async fn load_todos() -> Result<Vec<Todo>, ServerFnError> {
    use crafting::dbal;
    use leptos_actix::extract;

    let db = extract::<DbConnection>().await?;
    dbal::get_todos(&db).await.map_err(ServerFnError::new)
}

#[server]
pub async fn save_todo(todo: Todo, action: String) -> Result<(), ServerFnError> {
    use crafting::dbal;
    use leptos_actix::extract;

    let db = extract::<DbConnection>().await?;
    if action == "create" {
        dbal::create_todo(todo, &db)
            .await
            .map_err(ServerFnError::new)
            .map(|_| ())
    } else if action == "update" {
        dbal::update_todo(todo, &db)
            .await
            .map_err(ServerFnError::new)
            .map(|_| ())
    } else {
        Ok(())
    }
}

#[server]
pub async fn delete_todo(id: i32) -> Result<(), ServerFnError> {
    use crafting::dbal;
    use leptos_actix::extract;

    let db = extract::<DbConnection>().await?;
    dbal::delete_todo(id, &db)
        .await
        .map_err(ServerFnError::new)
        .map(|_| ())
}

#[component]
pub fn App() -> impl IntoView {
    provide_meta_context();

    view! {
        <Stylesheet id="leptos" href="/pkg/crafting.css" />
        <Link href="/assets/favicon.svg" rel="icon" type_="image/svg+xml" />
        <Link href="/assets/favicon.svg" rel="mask-icon" />

        <Meta content="width=device-width, initial-scale=1" name="viewport" />
        <Title text="Crafting Projekte" />

        <Router>
            <Routes fallback=|| AppPage>
                <Route path=path!("") view=AppPage />
                <Route path=path!("/*any") view=AppPage />
            </Routes>
        </Router>
    }
}

#[component]
fn AppPage() -> impl IntoView {
    let todos_resource = Resource::new(|| (), |_| async move { load_todos().await });
    let selected_todo = RwSignal::new(Todo::default());

    let dialog_ref = NodeRef::new();

    let create_action = ServerAction::<SaveTodo>::new();
    let delete_action = ServerAction::<DeleteTodo>::new();

    Effect::new(move |_| {
        if create_action.value().get().is_some_and(|res| res.is_ok()) {
            todos_resource.refetch();
        }
    });

    Effect::new(move |_| {
        if delete_action.value().get().is_some_and(|res| res.is_ok()) {
            if let Some(dialog_ref) = dialog_ref.get() as Option<HtmlDialogElement> {
                dialog_ref.close();
            }
            todos_resource.refetch();
        }
    });

    view! {
        <Transition>
            <main class="page">
                <h1 class="header">Projekte</h1>
                <aside class="sidebar">
                    {move || {
                        Suspend::new(async move {
                            todos_resource
                                .await
                                .ok()
                                .map(|todos| {
                                    if selected_todo.get().id == 0 && !selected_todo.get().is_new {
                                        selected_todo
                                            .set(
                                                todos
                                                    .first()
                                                    .cloned()
                                                    .unwrap_or(Todo {
                                                        is_new: true,
                                                        ..Default::default()
                                                    }),
                                            );
                                    }
                                    let (open_todos, closed_todos): (Vec<_>, Vec<_>) = todos
                                        .into_iter()
                                        .partition(|todo| !todo.is_done);
                                    view! {
                                        <For each=move || open_todos.clone() key=move |todo| todo.id let(todo)>
                                            <a class="sidebar__item" class:is--selected=selected_todo.get().id == todo.id on:click=move |_| selected_todo.set(todo.clone())>{todo.title.clone()}</a>
                                        </For>
                                        <a class="sidebar__item" class:is--selected=selected_todo.get().is_new on:click=move |_| selected_todo.set(Todo {
                                            is_new: true,
                                            ..Default::default()
                                        })>
                                            Neues Projekt
                                        </a>
                                        <hr class="sidebar__separator"/>
                                        <For each=move || closed_todos.clone() key=move |todo| todo.id let(todo)>
                                            <a class="sidebar__item is--done" class:is--selected=selected_todo.get().id == todo.id on:click=move |_| selected_todo.set(todo.clone())>{todo.title.clone()}</a>
                                        </For>
                                    }
                                })
                        })
                    }}
                </aside>
                <div class="content">
                    <h3>
                        <Show
                            when=move || !selected_todo.read().is_new
                            fallback=move || view! { Neues Projekt }
                        >
                            {move || selected_todo.read().title.clone()}
                        </Show>
                    </h3>
                    <ActionForm attr:class="form" action=create_action>
                        <Show
                            when=move || selected_todo.read().is_new
                            fallback=move || {
                                view! {
                                    <input
                                        type="hidden"
                                        name="todo[id]"
                                        prop:value=selected_todo.read().id
                                    />
                                    <input type="hidden" name="todo[title]" prop:value="" />
                                }
                            }
                        >
                            <label class="label" for="title">
                                Name
                            </label>
                            <input
                                class="input"
                                required
                                name="todo[title]"
                                id="title"
                                type="text"
                                prop:value=move || selected_todo.read().title.clone()
                            />
                        </Show>
                        <label class="label" for="kind">
                            Art
                        </label>
                        <input
                            class="input"
                            name="todo[kind]"
                            id="kind"
                            type="text"
                            prop:value=move || selected_todo.read().kind.clone()
                        />
                        <label class="label" for="status">
                            Status
                        </label>
                        <input
                            class="input"
                            name="todo[status]"
                            id="status"
                            type="text"
                            prop:value=move || selected_todo.read().status.clone()
                        />
                        <label class="label" for="material">
                            Material
                        </label>
                        <input
                            class="input"
                            name="todo[material]"
                            id="material"
                            type="text"
                            prop:value=move || selected_todo.read().material.clone()
                        />
                        <div class="input">
                            <input
                                class="checkbox"
                                id="isDone"
                                name="todo[is_done]"
                                type="checkbox"
                                value="true"
                                prop:checked=move || selected_todo.read().is_done
                            />
                            <label for="isDone">Erledigt</label>
                        </div>
                        <div class="input">
                            <input
                                class="checkbox"
                                id="isPartsMissing"
                                name="todo[is_parts_missing]"
                                type="checkbox"
                                value="true"
                                prop:checked=move || selected_todo.read().is_parts_missing
                            />
                            <label for="isPartsMissing">Teile Fehlen</label>
                        </div>
                        <label class="label" for="notes">
                            Notizen
                        </label>
                        <textarea
                            class="input"
                            name="todo[notes]"
                            id="notes"
                            rows="10"
                            prop:value=move || selected_todo.read().notes.clone()
                        ></textarea>
                        <div class="button-row">
                            <Show when=move || !selected_todo.read().is_new>
                                <button
                                    class="button is--negative"
                                    type="button"
                                    onclick="confirmDelete.showModal()"
                                >
                                    <svg viewBox="0 0 24 24" class="lucide">
                                        <path d="M3 6h18" />
                                        <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
                                        <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
                                    </svg>
                                    Löschen
                                </button>
                            </Show>
                            <button
                                type="submit"
                                class="button"
                                name="action"
                                value=move || {
                                    selected_todo
                                        .read()
                                        .is_new
                                        .then_some("create")
                                        .unwrap_or("update")
                                }
                            >
                                <svg viewBox="0 0 24 24" class="lucide">
                                    <path d="M15.2 3a2 2 0 0 1 1.4.6l3.8 3.8a2 2 0 0 1 .6 1.4V19a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2z" />
                                    <path d="M17 21v-7a1 1 0 0 0-1-1H8a1 1 0 0 0-1 1v7" />
                                    <path d="M7 3v4a1 1 0 0 0 1 1h7" />
                                </svg>
                                Speichern
                            </button>
                        </div>
                    </ActionForm>
                </div>
            </main>
            <dialog class="dialog" id="confirmDelete" node_ref=dialog_ref>
                <ActionForm action=delete_action>
                    <Show when=move || !selected_todo.read().is_new>
                        <input type="hidden" name="id" prop:value=move || selected_todo.read().id />
                    </Show>
                    <span class="dialog__title">
                        {move || format!("Projekt {} löschen?", selected_todo.read().title)}
                    </span>
                    <p class="dialog__content">
                        {move || {
                            format!(
                                "Möchtest du das Projekt {} wirklich löschen?",
                                selected_todo.read().title,
                            )
                        }} <Show when=move || selected_todo.read().is_done>
                            <br />
                            Das Projekt ist noch nicht erledigt bist du sicher?
                        </Show>
                    </p>
                    <div class="dialog__buttons">
                        <button
                            class="button is--dialog"
                            type="button"
                            onclick="confirmDelete.close()"
                        >
                            <svg viewBox="0 0 24 24" class="lucide">
                                <path d="M18 6 6 18" />
                                <path d="m6 6 12 12" />
                            </svg>
                            Nicht löschen
                        </button>
                        <button class="button is--dialog is--negative" type="submit">
                            <svg viewBox="0 0 24 24" class="lucide">
                                <path d="M3 6h18" />
                                <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
                                <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
                            </svg>
                            Projekt löschen
                        </button>
                    </div>
                </ActionForm>
            </dialog>
        </Transition>
    }
}
