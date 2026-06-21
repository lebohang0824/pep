# Pep Todo Application Example
#
# Demonstrates every Pep language feature:
# app, meta, action, feature, param (string/integer/enum),
# required/default, rules (empty/not), functions, events
# (conditional and unconditional), and comments.

app TodoApp
    meta
        name: "Todo Application"
        description: "A task management application"
        version: "1.0.0"
        backend: "laravel:12"
        frontend: "vue"
        database: "sqlite"
    end

    # Action with meta and typed parameters
    action SendCompletionEmail
        meta
            description: "Sends a completion notification"
        end

        params
            param TodoId
                type: integer
                required: true
            end

            param Email
                type: string
                required: true
            end
        end
    end

    # Action with an enum parameter type
    action NotifyAdmin
        meta
            description: "Notifies an administrator"
        end

        params
            param Priority
                type: enum{"LOW","NORMAL","HIGH"}
                required: true
            end
        end
    end

    # Feature demonstrating all sub-blocks
    feature CreateTodo
        meta
            description: "Creates a new todo item"
        end

        # A reusable validation function with a comparison
        function isStatusDone(s)
            return s == "DONE"
        end

        params
            param Title
                type: string
                required: true
            end

            param Description
                type: string
                required: false
                default: ""
            end

            param Priority
                type: integer
                required: false
                default: 1
            end

            param Status
                type: enum{"QUEUE","IN_PROGRESS","DONE"}
                required: false
                default: "QUEUE"
            end
        end

        rules
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

    feature UpdateTodo
        meta
            description: "Updates an existing todo item"
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

            param Assignee
                type: string
                required: false
                default: "unassigned"
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
            description: "Deletes a todo item"
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

        # Unconditional event — always fires
        events
            trigger NotifyAdmin
        end
    end

    feature GetTodo
        meta
            description: "Retrieves a single todo item"
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
            description: "Lists all todo items with optional filtering"
        end

        params
            param StatusFilter
                type: enum{"ALL","QUEUE","IN_PROGRESS","DONE"}
                required: false
                default: "ALL"
            end

            param PageSize
                type: integer
                required: false
                default: 20
            end
        end
    end

    # Feature demonstrating a function-based rule with not
    # and a function-call event condition
    feature MarkTodoDone
        meta
            description: "Marks a todo item as done"
        end

        function isHighPriority(p)
            return p >= 5
        end

        params
            param Id
                type: integer
                required: true
            end

            param Priority
                type: integer
                required: false
                default: 1
            end
        end

        rules
            rule IdRequired
                if empty(Id)
                    reject
            end

            # Uses the isHighPriority function — rejects if false
            rule MustBeHighPriority
                if not isHighPriority(Priority)
                    reject
            end
        end

        events
            if isHighPriority(Priority)
                trigger NotifyAdmin
            end
        end
    end
end
