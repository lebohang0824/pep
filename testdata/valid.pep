app TodoApp
    meta
        name: "Todo Application"
        description: "A simple task management application."
        version: "1.0.0"
        database: "sqlite"
        backend: "laravel:12"
        frontend: "vue"
    end

    action SendCompletionEmail
        params
            param TodoId
                type: integer
                required: true
            end
        end
    end

    feature CreateTodo
        meta
            description: "Creates a new todo item."
        end

        params
            param Title
                type: string
                required: true
            end
            param Description
                type: string
                required: false
            end
            param Status
                type: enum{"QUEUE","IN_PROGRESS","DONE"}
                default: "QUEUE"
            end
        end

        rules
            rule TitleRequired
                if empty(Title)
                    reject
            end
        end
    end

    feature UpdateTodo
        meta
            description: "Updates an existing todo item."
        end

        params
            param Id
                type: integer
                required: true
            end
            param Title
                type: string
                required: true
            end
            param Description
                type: string
                required: false
            end
            param Status
                type: enum{"QUEUE","IN_PROGRESS","DONE"}
                required: true
            end
        end

        rules
            rule IdRequired
                if empty(Id)
                    reject
            end
            rule TitleRequired
                if empty(Title)
                    reject
            end
        end

        events
            if Status == "DONE"
                trigger SendCompletionEmail
            end
        end
    end

    feature DeleteTodo
        meta
            description: "Deletes an existing todo item."
        end

        params
            param Id
                type: integer
                required: true
            end
        end

        rules
            rule IdRequired
                if empty(Id)
                    reject
            end
        end
    end

    feature GetTodo
        meta
            description: "Retrieves a single todo item."
        end

        params
            param Id
                type: integer
                required: true
            end
        end

        rules
            rule IdRequired
                if empty(Id)
                    reject
            end
        end
    end

    feature ListTodos
        meta
            description: "Returns all todo items."
        end
    end
end
