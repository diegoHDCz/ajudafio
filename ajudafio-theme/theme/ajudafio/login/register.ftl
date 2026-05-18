<#import "template.ftl" as layout>
<@layout.registrationLayout displayInfo=false displayWide=false; section>
    <#if section = "form">
        <div class="af-logo-wrap">
            <img src="${url.resourcesPath}/img/logo.png" alt="Ajuda Fio" class="af-form-logo" />
        </div>
        <h1 class="af-form-title">Criar conta</h1>
        <p class="af-form-sub">Preencha seus dados para se cadastrar</p>

        <form id="kc-register-form" action="${url.registrationAction}" method="post">
            <div class="af-field-row">
                <div class="af-field">
                    <label for="firstName">Nome</label>
                    <input type="text" id="firstName" name="firstName"
                        value="${(register.formData.firstName!'')}" placeholder="Nome" />
                    <#if messagesPerField.existsError('firstName')>
                        <span class="af-error">${kcSanitize(messagesPerField.get('firstName'))?no_esc}</span>
                    </#if>
                </div>
                <div class="af-field">
                    <label for="lastName">Sobrenome</label>
                    <input type="text" id="lastName" name="lastName"
                        value="${(register.formData.lastName!'')}" placeholder="Sobrenome" />
                    <#if messagesPerField.existsError('lastName')>
                        <span class="af-error">${kcSanitize(messagesPerField.get('lastName'))?no_esc}</span>
                    </#if>
                </div>
            </div>

            <div class="af-field">
                <label for="email">E-mail</label>
                <input type="text" id="email" name="email"
                    value="${(register.formData.email!'')}" autocomplete="email" placeholder="seu@email.com" />
                <#if messagesPerField.existsError('email')>
                    <span class="af-error">${kcSanitize(messagesPerField.get('email'))?no_esc}</span>
                </#if>
            </div>

            <#if !realm.registrationEmailAsUsername>
                <div class="af-field">
                    <label for="username">Usuário</label>
                    <input type="text" id="username" name="username"
                        value="${(register.formData.username!'')}" autocomplete="username" placeholder="nome de usuário" />
                    <#if messagesPerField.existsError('username')>
                        <span class="af-error">${kcSanitize(messagesPerField.get('username'))?no_esc}</span>
                    </#if>
                </div>
            </#if>

            <#if passwordRequired??>
                <div class="af-field">
                    <label for="password">Senha</label>
                    <input type="password" id="password" name="password"
                        autocomplete="new-password" placeholder="••••••••" />
                    <#if messagesPerField.existsError('password')>
                        <span class="af-error">${kcSanitize(messagesPerField.get('password'))?no_esc}</span>
                    </#if>
                </div>
                <div class="af-field">
                    <label for="password-confirm">Confirmar senha</label>
                    <input type="password" id="password-confirm" name="password-confirm"
                        autocomplete="new-password" placeholder="••••••••" />
                    <#if messagesPerField.existsError('password-confirm')>
                        <span class="af-error">${kcSanitize(messagesPerField.get('password-confirm'))?no_esc}</span>
                    </#if>
                </div>
            </#if>

            <button type="submit" class="af-btn">Cadastrar</button>
        </form>

        <div class="af-register-link">
            <span>Já tem uma conta?</span>
            <a href="${url.loginUrl}">Entrar</a>
        </div>
    </#if>
</@layout.registrationLayout>
