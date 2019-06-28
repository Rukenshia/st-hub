package sthub
{
   import lesta.components.LoginData;
   import lesta.constants.ComponentClass;
   import lesta.data.GameDelegate;
   import lesta.datahub.DataHub;
   import lesta.structs.DockData;
   import lesta.structs.ServerState;
   import lesta.unbound.core.UbController;
   import lesta.unbound.core.UbScope;
   import lesta.unbound.expression.IUbExpression;
   import lesta.utils.GameInfoHolder;
   
   public class Controller extends UbController
   {
      
      public function Controller()
      {
         super();
		 
		 scope.test = "Hallo";
		 scope.onEvent('test' + UbScope.CHANGED,null,UbScope.EVENT_DIRECTION_NONE);
      }
      
      override public function init(param1:Vector.<IUbExpression>) : void
      {
         super.init(param1);
      }
   }
}
